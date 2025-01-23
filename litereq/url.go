package litereq

import (
	"net/url"
	"sort"
	"strings"
)

type urlBuilder struct {
	baseurl                       string
	scheme, host                  string
	paths                         []string
	params, footParams, godParams []multimap
}

// ------------------ ub -----------------------

func (ub *urlBuilder) BaseURL(baseurl string) {
	ub.baseurl = baseurl
}

func (ub *urlBuilder) Scheme(scheme string) {
	ub.scheme = scheme
}

func (ub *urlBuilder) Host(host string) {
	ub.host = host

}

func (ub *urlBuilder) Path(path string) {
	ub.paths = append(ub.paths, path)
}

func (ub *urlBuilder) Param(key string, values ...string) {
	ub.params = append(ub.params, multimap{key, values})
}

func (ub *urlBuilder) GodParam(key string, values ...string) {
	ub.godParams = append(ub.godParams, multimap{key, values})
}

func (ub *urlBuilder) FootParam(key string, values ...string) {
	ub.footParams = append(ub.footParams, multimap{key, values})
}

func (ub *urlBuilder) Clone() *urlBuilder {
	ub2 := *ub
	Clip(&ub2.paths)
	Clip(&ub2.params)
	return &ub2
}

func (ub *urlBuilder) URL() (u *url.URL, err error) {
	u, err = url.Parse(ub.baseurl)
	if err != nil {
		return new(url.URL), err
	}
	u.Scheme = Or(
		ub.scheme,
		u.Scheme,
		"https",
	)
	u.Host = Or(ub.host, u.Host)
	for _, p := range ub.paths {
		u.Path = u.ResolveReference(&url.URL{Path: p}).Path
	}

	q := u.Query()
	if len(ub.params) > 0 {
		for _, kv := range ub.params {
			q[kv.key] = kv.values
		}
	}
	f := make(url.Values)
	if len(ub.footParams) > 0 {
		for _, kv := range ub.footParams {
			f[kv.key] = kv.values
		}
	}
	u.RawQuery = combineQuery(encode(q), encode(f))
	// Reparsing, in case the path rewriting broke the URL
	u, err = url.Parse(u.String())
	if err != nil {
		return new(url.URL), err
	}
	return u, nil
}

func encode(v url.Values) string {
	v.Encode()
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := k
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(queryEscape(keyEscaped))
			buf.WriteByte('=')
			buf.WriteString(queryEscape(v))
		}
	}
	return buf.String()
}
func queryEscape(v string) string {
	v = url.QueryEscape(v)
	//v = strings.ReplaceAll(v, "%2C", ",")
	v = strings.ReplaceAll(v, "+", "%20")
	v = strings.ReplaceAll(v, "%2A", "*")
	return v
}

func combineQuery(s ...string) string {
	q := strings.Builder{}
	if len(s) > 0 {
		for i, qs := range s {
			q.WriteString(qs)
			if i < len(s) {
				q.WriteString("&")
			}
		}
	}
	return strings.Trim(q.String(), "&")
}
