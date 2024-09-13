package mock

import (
	"context"
	"net/http"
)

var CancelForce context.CancelFunc

type ServerMock struct{}

func (s *ServerMock) ListenAndServe() error {
	CancelForce()
	return http.ErrServerClosed
}

func (ServerMock) Shutdown(ctx context.Context) error { return nil }
