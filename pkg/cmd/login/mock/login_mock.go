package mock

import "context"

type ServerMock struct{}

func (ServerMock) ListenAndServe() error { return nil }
func (ServerMock) Shutdown(ctx context.Context) error
