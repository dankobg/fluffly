package httpserver

import (
	"net"
	"net/http"
	"time"
)

// New returns the new [http.Server]
func New(opts ...ServerOption) *http.Server {
	sopts := &serverOpts{
		readTimeout:       10 * time.Second,
		readHeaderTimeout: 4 * time.Second,
		writeTimeout:      10 * time.Second,
		idleTimeout:       60 * time.Second,
	}

	for _, o := range opts {
		o.apply(sopts)
	}

	srv := &http.Server{
		Addr:                         sopts.addr,
		Handler:                      sopts.handler,
		DisableGeneralOptionsHandler: sopts.disableGeneralOptionsHandler,
		TLSConfig:                    sopts.tLSConfig,
		ReadTimeout:                  sopts.readTimeout,
		ReadHeaderTimeout:            sopts.readHeaderTimeout,
		WriteTimeout:                 sopts.writeTimeout,
		IdleTimeout:                  sopts.idleTimeout,
		MaxHeaderBytes:               sopts.maxHeaderBytes,
		TLSNextProto:                 sopts.tLSNextProto,
		ConnState:                    sopts.connState,
		ErrorLog:                     sopts.errorLog,
		BaseContext:                  sopts.baseContext,
		ConnContext:                  sopts.connContext,
	}

	return srv
}

func NewHttpClient() *http.Client {
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			MaxConnsPerHost:       10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return c
}
