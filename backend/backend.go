package backend

import (
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

func InitBackends(urls []string) ([]*Backend, error) {
	var result []*Backend
	for _, raw := range urls {
		parsed, err := url.Parse(raw)
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(parsed)
		result = append(result, &Backend{
			URL:          parsed,
			Alive:        true,
			ReverseProxy: proxy,
		})
	}
	return result, nil
}
