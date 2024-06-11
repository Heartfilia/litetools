package litereq

import "sync"

type Context struct {
	ctxMap map[string]any
	lock   *sync.RWMutex
}

func NewContext() *Context {
	return &Context{
		ctxMap: make(map[string]any),
		lock:   &sync.RWMutex{},
	}
}

func (c *Context) Put(key string, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ctxMap[key] = value
}

func (c *Context) Get(key string) any {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.ctxMap[key]; ok {
		return value
	}
	return nil
}

func (c *Context) GetString(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.ctxMap[key]; ok {
		return value.(string)
	}
	return ""
}

func (c *Context) GetBool(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.ctxMap[key]; ok {
		return value.(bool)
	}
	return false
}

func (c *Context) GetInt(key string) int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.ctxMap[key]; ok {
		switch value.(type) {
		case int64:
			return int(value.(int64))
		case int32:
			return int(value.(int32))
		case int:
			return value.(int)
		}
		return 0
	}
	return 0
}

func (c *Context) Items() []any {
	c.lock.RLock()
	defer c.lock.RUnlock()
	ret := make([]any, 0, len(c.ctxMap))
	for k, v := range c.ctxMap {
		ret = append(ret, map[string]any{k: v})
	}
	return ret
}

func (c *Context) Keys() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	ret := make([]string, 0, len(c.ctxMap))
	for k, _ := range c.ctxMap {
		ret = append(ret, k)
	}
	return ret
}

func (c *Context) Values() []any {
	c.lock.RLock()
	defer c.lock.RUnlock()
	ret := make([]any, 0, len(c.ctxMap))
	for _, v := range c.ctxMap {
		ret = append(ret, v)
	}
	return ret
}
