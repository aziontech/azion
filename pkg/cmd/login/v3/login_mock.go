package login

import (
	"context"
	"net/http"
)

type ServerMock struct {
	Cancel               bool
	ErrorListenAndServer error
	ErrorShutdown        error
}

func (s *ServerMock) ListenAndServe() error {
	if s.Cancel {
		shutdownContext()
	}
	if s.ErrorListenAndServer != nil {
		return s.ErrorListenAndServer
	}
	return http.ErrServerClosed
}

func (s *ServerMock) Shutdown(ctx context.Context) error {
	if s.ErrorShutdown != nil {
		return s.ErrorShutdown
	}
	return nil
}
