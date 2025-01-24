package litereq

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// 借鉴 https://github.com/earthboundkid/requests

// Transport is an alias of http.RoundTripper for documentation purposes.
type Transport = http.RoundTripper

// RoundTripFunc is an adaptor to use a function as an http.RoundTripper.
type RoundTripFunc func(req *http.Request) (res *http.Response, err error)

// RoundTrip implements http.RoundTripper.
func (rtf RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rtf(r)
}

var _ Transport = RoundTripFunc(nil)

var (
	tLock            sync.RWMutex
	transports       = make(map[string]*http.Transport)
	defaultTransport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		MaxIdleConnsPerHost:   300,
		TLSClientConfig:       nil,
	}
)

// ReplayString returns an http.RoundTripper that always responds with a
// request built from rawResponse. It is intended for use in one-off tests.
func ReplayString(rawResponse string) Transport {
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r := bufio.NewReader(strings.NewReader(rawResponse))
		res, err = http.ReadResponse(r, req)
		return
	})
}

// UserAgentTransport returns a wrapped http.RoundTripper that sets the User-Agent header.
func UserAgentTransport(rt http.RoundTripper, s string) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r2 := *req
		r2.Header = r2.Header.Clone()
		r2.Header.Set("User-Agent", s)
		return rt.RoundTrip(&r2)
	})
}

// PermitURLTransport returns a wrapped http.RoundTripper that rejects any url whose URL doesn't match the provided regular expression string.
//
// PermitURLTransport will panic if the regexp does not compile.
func PermitURLTransport(rt http.RoundTripper, regex string) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	re := regexp.MustCompile(regex)
	reErr := fmt.Errorf("requested URL not permitted by regexp: %s", regex)
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		if u := req.URL.String(); !re.MatchString(u) {
			return nil, reErr
		}
		return rt.RoundTrip(req)
	})
}

// LogTransport returns a wrapped http.RoundTripper
// that calls fn with details when a response has finished.
// A response is considered finished
// when the wrapper http.RoundTripper returns an error
// or the Response.Body is closed,
// whichever comes first.
// To simplify logging code,
// a nil *http.Response is replaced with a new http.Response.
func LogTransport(rt http.RoundTripper, fn func(req *http.Request, res *http.Response, err error, duration time.Duration)) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		start := time.Now()
		res, err = rt.RoundTrip(req)
		if err != nil {
			res2 := res
			if res == nil {
				res2 = new(http.Response)
			}
			fn(req, res2, err, time.Since(start))
			return
		}

		res.Body = closeLogger{res.Body, func() {
			fn(req, res, err, time.Since(start))
		}}
		return
	})
}

type closeLogger struct {
	io.ReadCloser
	fn func()
}

func (cl closeLogger) Close() error {
	cl.fn()
	return cl.ReadCloser.Close()
}

// DoerTransport converts a Doer into a Transport.
// It exists for compatibility with other libraries.
// A Doer is an interface with a Do method.
// Users should prefer Transport,
// because Do is the interface of http.Client
// which has higher level concerns.
func DoerTransport(cl interface {
	Do(req *http.Request) (*http.Response, error)
}) Transport {
	return RoundTripFunc(cl.Do)
}

func getTransport(url string) *http.Transport {
	tLock.RLock()
	defer tLock.RUnlock()
	t, ok := transports[url]
	if ok {
		return t
	}
	return nil
}

func createTransport(getter ProxyGetter, h2 bool) *http.Transport {
	if getter == nil {
		return defaultTransport
	}
	proxy := getter()
	if proxy == nil {
		return defaultTransport
	}
	url := proxy.URL()
	if url == nil {
		return defaultTransport
	}
	t := getTransport(proxy.String())
	if t != nil {
		return t
	}
	tLock.Lock()
	defer tLock.Unlock()
	t = &http.Transport{
		Proxy:                 http.ProxyURL(url),
		ForceAttemptHTTP2:     h2,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		MaxIdleConnsPerHost:   1000,
		TLSClientConfig:       nil,
	}
	transports[proxy.String()] = t
	time.AfterFunc(proxy.Expired(), func() {
		tLock.Lock()
		defer tLock.Unlock()
		delete(transports, proxy.String())
	})
	return t
}
