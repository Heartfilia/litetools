package litepage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"slices"
	"sync"
	"sync/atomic"
)

type cdpRequest struct {
	ID     int64  `json:"id"`
	Method string `json:"method"`
	Params any    `json:"params,omitempty"`
}

type cdpResponse struct {
	ID     int64           `json:"id,omitempty"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *cdpError       `json:"error,omitempty"`
}

type cdpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *cdpError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("cdp error %d: %s", e.Code, e.Message)
}

type cdpConn struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
	waitMu  sync.Mutex
	waiters map[int64]chan cdpResponse
	subMu   sync.RWMutex
	subs    map[string][]chan json.RawMessage
	nextID  atomic.Int64
	closed  atomic.Bool
}

func newCDPConn(wsURL string) (*cdpConn, error) {
	header := map[string][]string{
		"Origin": {originFromWSURL(wsURL)},
	}
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		return nil, err
	}
	c := &cdpConn{
		conn:    conn,
		waiters: make(map[int64]chan cdpResponse),
		subs:    make(map[string][]chan json.RawMessage),
	}
	go c.readLoop()
	return c, nil
}

func (c *cdpConn) readLoop() {
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			c.failAll(err)
			return
		}

		var resp cdpResponse
		if err := json.Unmarshal(raw, &resp); err != nil {
			continue
		}
		if resp.ID == 0 {
			if resp.Method != "" {
				c.dispatchEvent(resp.Method, resp.Params)
			}
			continue
		}

		c.waitMu.Lock()
		ch, ok := c.waiters[resp.ID]
		if ok {
			delete(c.waiters, resp.ID)
		}
		c.waitMu.Unlock()
		if ok {
			ch <- resp
		}
	}
}

func (c *cdpConn) dispatchEvent(method string, payload json.RawMessage) {
	c.subMu.RLock()
	subs := append([]chan json.RawMessage(nil), c.subs[method]...)
	c.subMu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- payload:
		default:
		}
	}
}

func (c *cdpConn) failAll(err error) {
	if !c.closed.CompareAndSwap(false, true) {
		return
	}
	c.waitMu.Lock()
	defer c.waitMu.Unlock()
	for id, ch := range c.waiters {
		delete(c.waiters, id)
		ch <- cdpResponse{ID: id, Error: &cdpError{Code: -1, Message: err.Error()}}
	}
	c.subMu.Lock()
	defer c.subMu.Unlock()
	for method, subs := range c.subs {
		for _, ch := range subs {
			close(ch)
		}
		delete(c.subs, method)
	}
}

func (c *cdpConn) close() error {
	if !c.closed.CompareAndSwap(false, true) {
		return nil
	}
	c.subMu.Lock()
	for method, subs := range c.subs {
		for _, ch := range subs {
			close(ch)
		}
		delete(c.subs, method)
	}
	c.subMu.Unlock()
	return c.conn.Close()
}

func (c *cdpConn) call(method string, params any, out any) error {
	if c == nil {
		return errors.New("nil cdp connection")
	}
	id := c.nextID.Add(1)
	req := cdpRequest{ID: id, Method: method, Params: params}
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ch := make(chan cdpResponse, 1)
	c.waitMu.Lock()
	c.waiters[id] = ch
	c.waitMu.Unlock()

	c.writeMu.Lock()
	err = c.conn.WriteMessage(websocket.TextMessage, payload)
	c.writeMu.Unlock()
	if err != nil {
		c.waitMu.Lock()
		delete(c.waiters, id)
		c.waitMu.Unlock()
		return err
	}

	resp := <-ch
	if resp.Error != nil {
		return resp.Error
	}
	if out == nil || len(resp.Result) == 0 {
		return nil
	}
	return json.Unmarshal(resp.Result, out)
}

func (c *cdpConn) subscribe(method string, buffer int) (<-chan json.RawMessage, func()) {
	if buffer <= 0 {
		buffer = 1
	}
	ch := make(chan json.RawMessage, buffer)
	c.subMu.Lock()
	c.subs[method] = append(c.subs[method], ch)
	c.subMu.Unlock()

	cancel := func() {
		c.subMu.Lock()
		defer c.subMu.Unlock()
		subs := c.subs[method]
		idx := slices.Index(subs, chan json.RawMessage(ch))
		if idx >= 0 {
			subs = append(subs[:idx], subs[idx+1:]...)
			if len(subs) == 0 {
				delete(c.subs, method)
			} else {
				c.subs[method] = subs
			}
		}
		close(ch)
	}
	return ch, cancel
}
