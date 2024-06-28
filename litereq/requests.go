package litereq

import (
	"errors"
	"github.com/Heartfilia/litetools/litereq/opt"
	"github.com/Heartfilia/litetools/litestr"
	"io"
	"log"
	netHTTP "net/http"
	netURL "net/url"
	"strings"
)

type Request struct {
	URL     *netURL.URL
	Ctx     *Context
	Options *opt.Option
}

func (s *Session) sendRequest(url string, o *opt.Option) *Response {
	response := NewResponse()
	suc := false
	for r := 0; r < s.maxRetry; r++ {
		if s.http2 {
			// 如果是http2模式下走这个地方 现在先不兼容
			log.Println(litestr.W(), "暂时还没兼容http2")
		} else {
			resp, respBytes, err := s.http1Request(url, o)
			if err != nil {
				if s.verbose {
					log.Println(litestr.E(), err)
				}
				response.err = err
				continue
			}
			response.Body = respBytes
			response.Text = string(respBytes)
			response.Headers = resp.Header
			response.StatusCode = resp.StatusCode
			response.Proto = resp.Proto
			response.Status = resp.Status
			response.ContentLength = int(resp.ContentLength)
			s.updateCookies(resp.Cookies()) // 保存cookie  >>> maybe 30x not success
			respCk := &opt.Cookie{}
			respCk.StoreCookies(resp.Cookies())
			response.Cookies = respCk
			response.err = nil
			suc = true
		}
		if s.verbose && response.Error() != nil {
			// 这里是在过程中遇到的报错打印出来
			log.Println(litestr.E(), "error:", response.Error())
		}
		if suc == true {
			break
		}
	}
	if suc == false && response.Error() == nil {
		// 如果失败的时候 并且没有失败的日志记录 那么补充一个错误提示
		response.err = errors.New("bad requests with this packages: help me fix it with debug")
	}

	return response
}

func (s *Session) http1Request(url string, o *opt.Option) (*netHTTP.Response, []byte, error) {
	var req *netHTTP.Request
	var err error
	baseNewContentType := ""
	switch o.GetMethod() {
	case "POST", "PUT", "DELETE", "PATCH":
		var body string
		if o.GetJson() != nil {
			body = string(o.GetJson())
			baseNewContentType = "application/json"
		} else {
			//dataInfo, typeInfo := o.GetData()
			//if typeInfo == "bytes" {
			//	baseNewContentType = "application/octet-stream"
			//} else if typeInfo == "form" {
			//	baseNewContentType = "application/x-www-form-urlencoded"
			//}

			return nil, nil, errors.New("not support now")
		}
		payload := strings.NewReader(body)
		req, err = netHTTP.NewRequest(o.GetMethod(), url, payload)
	case "OPTIONS", "GET", "HEAD", "TRACE":
		req, err = netHTTP.NewRequest(o.GetMethod(), url, nil)
	case "CONNECT":
		log.Panicln("暂时不支持 not support now")
	default:
		log.Panicf("not support your method: %s", o.GetMethod())
	}
	if err != nil {
		return nil, nil, err
	}

	if o.GetParams() != nil {
		req.URL.RawQuery = o.GetParams().Encode()
	}
	if baseNewContentType != "" {
		req.Header.Set("Content-Type", baseNewContentType)
	} // 先程序自动配置header类型，然后下面再参数补充
	s.setReqHeaders(req, o.GetHeaders())
	if o.GetCookieEnable() {
		s.setReqCookies(req, o.GetCookies())
	}
	s.setTimeout(o.GetTimeout())
	s.setProxy(o.GetProxy())
	if s.host != "" {
		req.Host = s.host
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	respByte, err := io.ReadAll(resp.Body)

	return resp, respByte, err
}
