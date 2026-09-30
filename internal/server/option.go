package server

import (
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/desulaidovich/gym/config"
)

type Option interface {
	apply(*Server) error
}

type handlerOption struct {
	handler http.Handler
}

func WithHandler(h http.Handler) Option {
	return handlerOption{
		handler: h,
	}
}

func (o handlerOption) apply(s *Server) error {
	if o.handler == nil {
		return fmt.Errorf("nil HTTP handler")
	}
	s.srv.Handler = o.handler
	return nil
}

type loggerOption struct {
	log *zap.Logger
}

func WithLogger(log *zap.Logger) Option {
	return loggerOption{
		log: log,
	}
}

func (o loggerOption) apply(s *Server) error {
	if o.log == nil {
		return fmt.Errorf("nil logger")
	}
	s.log = o.log
	return nil
}

type configOption struct {
	cfg *config.HTTP
}

func WithConfig(cfg *config.HTTP) Option {
	return configOption{
		cfg: cfg,
	}
}

func (o configOption) apply(s *Server) error {
	if o.cfg == nil {
		return fmt.Errorf("nil HTTP config")
	}

	s.addr = net.JoinHostPort(o.cfg.Host, o.cfg.Port)
	if o.cfg.Network != "" {
		s.network = o.cfg.Network
	}

	s.srv.ReadHeaderTimeout = o.cfg.ReadHeaderTimeout
	s.srv.ReadTimeout = o.cfg.ReadTimeout
	s.srv.WriteTimeout = o.cfg.WriteTimeout
	s.srv.IdleTimeout = o.cfg.IdleTimeout
	s.srv.MaxHeaderBytes = o.cfg.MaxHeaderBytes
	s.srv.MaxHeaderValueCount = o.cfg.MaxHeaderValueCount
	s.srv.DisableGeneralOptionsHandler = o.cfg.DisableGeneralOptionsHandler
	s.srv.DisableClientPriority = o.cfg.DisableClientPriority
	s.disableKeepAlives = o.cfg.DisableKeepAlives

	if o.cfg.ShutdownTimeout > 0 {
		s.shutdownTimeout = o.cfg.ShutdownTimeout
	}
	if o.cfg.Protocols != nil {
		var protocols http.Protocols
		protocols.SetHTTP1(o.cfg.Protocols.HTTP1)
		protocols.SetHTTP2(o.cfg.Protocols.HTTP2)
		protocols.SetUnencryptedHTTP2(o.cfg.Protocols.UnencryptedHTTP2)
		s.srv.Protocols = &protocols
	}
	if o.cfg.HTTP2 != nil {
		s.srv.HTTP2 = &http.HTTP2Config{
			MaxConcurrentStreams:          o.cfg.HTTP2.MaxConcurrentStreams,
			MaxDecoderHeaderTableSize:     o.cfg.HTTP2.MaxDecoderHeaderTableSize,
			MaxEncoderHeaderTableSize:     o.cfg.HTTP2.MaxEncoderHeaderTableSize,
			MaxReadFrameSize:              o.cfg.HTTP2.MaxReadFrameSize,
			MaxReceiveBufferPerConnection: o.cfg.HTTP2.MaxReceiveBufferPerConnection,
			MaxReceiveBufferPerStream:     o.cfg.HTTP2.MaxReceiveBufferPerStream,
			SendPingTimeout:               o.cfg.HTTP2.SendPingTimeout,
			PingTimeout:                   o.cfg.HTTP2.PingTimeout,
			WriteByteTimeout:              o.cfg.HTTP2.WriteByteTimeout,
		}
	}
	if o.cfg.TCPKeepAlive != nil {
		s.tcpKeepAlive = &net.KeepAliveConfig{
			Enable:   o.cfg.TCPKeepAlive.Enable,
			Idle:     o.cfg.TCPKeepAlive.Idle,
			Interval: o.cfg.TCPKeepAlive.Interval,
			Count:    o.cfg.TCPKeepAlive.Count,
		}
	}
	return nil
}
