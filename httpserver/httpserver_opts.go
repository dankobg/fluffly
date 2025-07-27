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
	tLSNextProto                 map[string]func(*http.Server, *tls.Conn, http.Handler)
	connState                    func(net.Conn, http.ConnState)
	errorLog                     *log.Logger
	baseContext                  func(net.Listener) context.Context
	connContext                  func(ctx context.Context, c net.Conn) context.Context
}

type ServerOption interface {
	apply(*serverOpts)
}

type ServerOptions []ServerOption

func (o ServerOptions) apply(s *serverOpts) {
	for _, opt := range o {
		opt.apply(s)
	}
}

type addrOpt string

func (o addrOpt) apply(s *serverOpts)   { s.addr = string(o) }
func WithAddr(addr string) ServerOption { return addrOpt(addr) }
func WithHostPort(host string, port int) ServerOption {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	return addrOpt(addr)
}

type handlerOpt struct{ h http.Handler }

func (o handlerOpt) apply(s *serverOpts)      { s.handler = o.h }
func WithHandler(h http.Handler) ServerOption { return handlerOpt{h: h} }

type readTimeoutOpt time.Duration

func (o readTimeoutOpt) apply(s *serverOpts)         { s.readTimeout = time.Duration(o) }
func WithReadTimeout(dur time.Duration) ServerOption { return readTimeoutOpt(dur) }

type readHeaderTimeoutOpt time.Duration

func (o readHeaderTimeoutOpt) apply(s *serverOpts)         { s.readHeaderTimeout = time.Duration(o) }
func WithReadHeaderTimeout(dur time.Duration) ServerOption { return readHeaderTimeoutOpt(dur) }

type writeTimeoutOpt time.Duration

func (o writeTimeoutOpt) apply(s *serverOpts)         { s.writeTimeout = time.Duration(o) }
func WithWriteTimeout(dur time.Duration) ServerOption { return writeTimeoutOpt(dur) }

type idleTimeoutOpt time.Duration

func (o idleTimeoutOpt) apply(s *serverOpts)         { s.idleTimeout = time.Duration(o) }
func WithIdleTimeout(dur time.Duration) ServerOption { return idleTimeoutOpt(dur) }

type maxHeaderBytesOpt int

func (o maxHeaderBytesOpt) apply(s *serverOpts) { s.maxHeaderBytes = int(o) }
func WithMaxHeaderBytes(bytes int) ServerOption { return maxHeaderBytesOpt(bytes) }

type disableGeneralOptionsHandlerOpt bool

func (o disableGeneralOptionsHandlerOpt) apply(s *serverOpts) {
	s.disableGeneralOptionsHandler = bool(o)
}
func WithDisableGeneralOptionsHandler(flag bool) ServerOption {
	return disableGeneralOptionsHandlerOpt(flag)
}

type errorLogOpt struct{ log *log.Logger }

func (o errorLogOpt) apply(s *serverOpts)       { s.errorLog = o.log }
func WithErrorLog(log *log.Logger) ServerOption { return errorLogOpt{log: log} }
func WithErrorSlog(sl *slog.Logger, level slog.Level) ServerOption {
	return errorLogOpt{log: slog.NewLogLogger(sl.Handler(), level)}
}

type tLSConfigOpt struct{ tls *tls.Config }

func (o tLSConfigOpt) apply(s *serverOpts)       { s.tLSConfig = o.tls }
func WithTLSConfig(tls *tls.Config) ServerOption { return tLSConfigOpt{tls: tls} }

type tLSNextProtoOpt map[string]func(*http.Server, *tls.Conn, http.Handler)

func (o tLSNextProtoOpt) apply(s *serverOpts) { s.tLSNextProto = o }
func WithTLSNextProto(proto map[string]func(*http.Server, *tls.Conn, http.Handler)) ServerOption {
	return tLSNextProtoOpt(proto)
}

type connStateOpt func(net.Conn, http.ConnState)

func (o connStateOpt) apply(s *serverOpts)                            { s.connState = o }
func WithConnState(state func(net.Conn, http.ConnState)) ServerOption { return connStateOpt(state) }

type baseContextOpt func(net.Listener) context.Context

func (o baseContextOpt) apply(s *serverOpts) { s.baseContext = o }
func WithBaseContext(basectx func(net.Listener) context.Context) ServerOption {
	return baseContextOpt(basectx)
}

type connContextOpt func(ctx context.Context, c net.Conn) context.Context

func (o connContextOpt) apply(s *serverOpts) { s.connContext = o }
func WithConnContext(connctx func(ctx context.Context, c net.Conn) context.Context) ServerOption {
	return connContextOpt(connctx)
}
