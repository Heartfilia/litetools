package litereq

import (
	"github.com/Heartfilia/litetools/litereq/opt"
	netURL "net/url"
)

type Request struct {
	URL     *netURL.URL
	Ctx     *Context
	Options *opt.Option
}
