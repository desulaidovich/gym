package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	DefaultShutdownTimeout = 15 * time.Second
	DefaultNetwork         = "tcp"
)

type Server struct {
	srv               *http.Server
	addr              string
	network           string
	shutdownTimeout   time.Duration
	disableKeepAlives bool
	tcpKeepAlive      *net.KeepAliveConfig
	log               *zap.Logger
}

func New(opts ...Option) (*Server, error) {
	s := &Server{
		srv:             new(http.Server),
		network:         DefaultNetwork,
		shutdownTimeout: DefaultShutdownTimeout,
	}

	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("server: nil option")
		}
		if err := opt.apply(s); err != nil {
			return nil, fmt.Errorf("server: %w", err)
		}
	}

	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	var lc net.ListenConfig
	if s.tcpKeepAlive != nil {
		lc.KeepAliveConfig = *s.tcpKeepAlive
	}

	ln, err := lc.Listen(ctx, s.network, s.addr)
	if err != nil {
		return fmt.Errorf("listen %q: %w", s.addr, err)
	}

	if s.disableKeepAlives {
		s.srv.SetKeepAlivesEnabled(false)
	}

	errCh := make(chan error, 1)
	go func() {
		s.log.Info("started",
			zap.String("addr", ln.Addr().String()),
			zap.String("network", s.network),
		)
		errCh <- s.srv.Serve(ln)
	}()

	select {
	case err = <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve %q: %w", s.addr, err)
		}
		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		if err = s.srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown %q: %w", s.addr, err)
		}
		if err = <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve %q: %w", s.addr, err)
		}

		s.log.Info("stopped", zap.String("addr", s.addr))
		return nil
	}
}
