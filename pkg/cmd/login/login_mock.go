package login

import (
	"context"
	"net/http"
)

type ServerMock struct{}

func (s *ServerMock) ListenAndServe() error {
	shutdownContext()
	return http.ErrServerClosed
}

func (ServerMock) Shutdown(ctx context.Context) error { return nil }
