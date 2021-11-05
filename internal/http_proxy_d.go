package internal

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

var (
	Title         = "http-proxy/0.1.0"
	ProxyID       = uuid.NewV4()
	ProxyIDHeader = "X-Proxy-ID"
)

type HttpProxyD struct {
	serv *http.Server
	conf *Config
}

func NewHttpProxyD(conf *Config) (HttpProxyD, error) {
	hp, err := NewHttpProxy(
		conf.Hosts,
		WithBalancerOfType(conf.BalancerType),
		WithBufferSize(conf.BufferSizeKb),
	)
	if err != nil {
		return HttpProxyD{}, err
	}

	hpd := HttpProxyD{
		serv: &http.Server{
			Addr:           conf.Host + ":" + strconv.Itoa(conf.Port),
			Handler:        hp,
			MaxHeaderBytes: conf.MaxHeaderKb,
		},
		conf: conf,
	}

	return hpd, nil
}

func (d HttpProxyD) Run() error {
	go d.serv.ListenAndServe()

	if err := d.checkHosts(); err != nil {
		d.serv.Shutdown(context.Background())
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	d.serv.Shutdown(context.Background())

	return nil
}

func (d HttpProxyD) checkHosts() error {
	for _, host := range d.conf.Hosts {
		req, err := http.NewRequest(http.MethodOptions, "http://"+host, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if resp.Header.Get(ProxyIDHeader) != "" {
			proxyID, err := uuid.FromString(resp.Header.Get(ProxyIDHeader))
			if err != nil {
				return err
			}

			if proxyID == ProxyID {
				return errors.New("distributed deadlock detected")
			}
		}
	}

	return nil
}
