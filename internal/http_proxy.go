package internal

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

type HttpProxy struct {
	hosts    []string
	buffer   []byte
	balancer Balancer
}

func NewHttpProxy(hosts []string, opts ...option) (*HttpProxy, error) {
	hp := &HttpProxy{
		hosts: hosts,
	}

	bufferSize := DefaultBufferSize

	for _, opt := range opts {
		if opt.err != nil {
			return nil, opt.err
		}
		if opt.balancer != nil {
			hp.balancer = opt.balancer
		}
		if opt.bufferSize != 0 {
			bufferSize = opt.bufferSize
		}
	}

	if hp.balancer == nil {
		hp.balancer = DefaultBalancer
	}

	hp.buffer = make([]byte, bufferSize)

	return hp, nil
}

func (hp HttpProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	err := hp.proxy(rw, req)
	if err != nil {
		rw.WriteHeader(502)
		rw.Write([]byte("Bad gateway"))
		return
	}
}

func (hp HttpProxy) proxy(rw http.ResponseWriter, req *http.Request) error {
	req = hp.newProxyRequest(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("bad gateway")
	}

	copyHeader(resp.Header, rw.Header())
	rw.Header().Add("Server", Title)
	rw.Header().Add(ProxyIDHeader, ProxyID.String())

	_, err = io.CopyBuffer(rw, resp.Body, hp.buffer)

	return err
}

func (hp HttpProxy) newProxyRequest(req *http.Request) *http.Request {
	var user *url.Userinfo
	if req.URL.User != nil {
		user = new(url.Userinfo)
		*user = *req.URL.User
	}

	return &http.Request{
		Method:        req.Method,
		Header:        req.Header,
		Form:          req.Form,
		PostForm:      req.PostForm,
		Body:          req.Body,
		MultipartForm: req.MultipartForm,
		URL: &url.URL{
			Scheme:      "http",
			Opaque:      req.URL.Opaque,
			User:        user,
			Host:        hp.balancer.Host(hp.hosts),
			Path:        req.URL.Path,
			RawPath:     req.URL.RawPath,
			ForceQuery:  req.URL.ForceQuery,
			RawQuery:    req.URL.RawQuery,
			Fragment:    req.URL.Fragment,
			RawFragment: req.URL.RawFragment,
		},
	}
}

func copyHeader(src http.Header, dst http.Header) {
	for hdr, value := range src {
		for _, subvalue := range value {
			dst.Add(hdr, subvalue)
		}
	}
}
