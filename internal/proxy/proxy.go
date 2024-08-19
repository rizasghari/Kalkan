package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/rizasghari/kalkan/internal/errs"
)

type Proxy struct {
	// ReverseProxy is an HTTP Handler that takes an incoming request and
	// sends it to another server, proxying the response back to the client.
	ReverseProxy *httputil.ReverseProxy
}

func New(target *url.URL) (*Proxy, error) {
	// Todo: Add more validations
	if target == nil || target.String() == "" {
		return nil, errs.ErrInvalidTargetUrl
	}
	rp := httputil.NewSingleHostReverseProxy(target)
	return &Proxy{ReverseProxy: rp}, nil
}

func (p *Proxy) ProxyRequestHandler(
	url *url.URL,
	endpoint string,
) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Kalkan: Request received at %s at %s\n", r.URL, time.Now().UTC())

		// Update the headers to allow for SSL redirection
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host

		// Trim reverseProxyRoutePrefix
		path := r.URL.Path
		r.URL.Path = strings.TrimLeft(path, endpoint)

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		fmt.Printf("Kalkan: Redirecting request to %s at %s\n", r.URL, time.Now().UTC())
		p.ReverseProxy.ServeHTTP(w, r)
	}
}
