package httpserver

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

type serverOpts struct {
	addr                         string
	handler                      http.Handler
	disableGeneralOptionsHandler bool
	tLSConfig                    *tls.Config
	readTimeout                  time.Duration
	readHeaderTimeout            time.Duration
	writeTimeout                 time.Duration
	idleTimeout                  time.Duration
	maxHeaderBytes               int
	maxHeaderValueCount          int
	tLSNextProto                 map[string]func(*http.Server, *tls.Conn, http.Handler)
	connState                    func(net.Conn, http.ConnState)
	errorLog                     *log.Logger
	baseContext                  func(net.Listener) context.Context
	connContext                  func(ctx context.Context, c net.Conn) context.Context
	http2                        *http.HTTP2Config
	protocols                    *http.Protocols
	disableClientPriority        bool
}

type ServerOption interface {
	apply(so *serverOpts)
}

type serverOptionFunc func(*serverOpts)

func (f serverOptionFunc) apply(s *serverOpts) {
	f(s)
}

type ServerOptions []ServerOption

func (o ServerOptions) apply(so *serverOpts) {
	for _, opt := range o {
		opt.apply(so)
	}
}

func WithAddr(addr string) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.addr = addr
	})
}

func WithHostPort(host string, port int) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		addr := net.JoinHostPort(host, strconv.Itoa(port))
		so.addr = addr
	})
}

func WithHandler(h http.Handler) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.handler = h
	})
}

func WithReadTimeout(dur time.Duration) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.readTimeout = dur
	})
}

func WithReadHeaderTimeout(dur time.Duration) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.readHeaderTimeout = dur
	})
}

func WithWriteTimeout(dur time.Duration) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.writeTimeout = dur
	})
}

func WithIdleTimeout(dur time.Duration) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.idleTimeout = dur
	})
}

func WithMaxHeaderBytes(n int) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.maxHeaderBytes = n
	})
}

func WithDisableGeneralOptionsHandler(flag bool) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.disableGeneralOptionsHandler = flag
	})
}

func WithErrorLog(log *log.Logger) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.errorLog = log
	})
}

func WithErrorSlog(sl *slog.Logger, level slog.Level) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.errorLog = slog.NewLogLogger(sl.Handler(), level)
	})
}

func WithTLSConfig(tls *tls.Config) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.tLSConfig = tls
	})
}

func WithConnState(state func(net.Conn, http.ConnState)) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.connState = state
	})
}

func WithBaseContext(f func(net.Listener) context.Context) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.baseContext = f
	})
}

func WithConnContext(f func(ctx context.Context, c net.Conn) context.Context) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.connContext = f
	})
}

func WithMaxHeaderValueCount(n int) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.maxHeaderValueCount = n
	})
}

func WithHTTP2(cfg *http.HTTP2Config) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.http2 = cfg
	})
}

func WithProtocols(protocols *http.Protocols) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.protocols = protocols
	})
}

func WithDisableClientPriority(flag bool) ServerOption {
	return serverOptionFunc(func(so *serverOpts) {
		so.disableClientPriority = flag
	})
}
