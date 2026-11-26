package version

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

func TestStates(t *testing.T) {
	tests := []struct {
		state           State
		known, settled  bool
		failed, done    bool
		buildFailsAfter bool
	}{
		{Draft, true, true, false, false, false},
		{Queued, true, false, false, false, false},
		{Building, true, false, false, false, false},
		{Ready, true, true, false, true, false},
		{Error, true, true, true, true, true},
		{Canceled, true, true, true, true, true},
		{Archiving, true, false, false, false, false},
		{Archived, true, true, true, true, true},
		{Deleting, true, false, false, false, false},
		{Deleted, true, true, true, true, true},
		{"promoting", false, false, false, false, false},
		{"", false, false, false, false, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			assert.Equal(t, tt.known, tt.state.Known(), "Known")
			assert.Equal(t, tt.settled, tt.state.Settled(), "Settled")
			assert.Equal(t, tt.failed, tt.state.Failed(), "Failed")

			done, err := BuildDone(tt.state)
			assert.Equal(t, tt.done, done, "BuildDone")
			assert.Equal(t, tt.buildFailsAfter, errors.Is(err, ErrBuildFailed), "BuildDone error")
		})
	}
}

func TestFailedImpliesSettled(t *testing.T) {
	for _, s := range []State{Draft, Queued, Building, Ready, Error, Canceled, Archiving, Archived, Deleting, Deleted} {
		if s.Failed() {
			assert.True(t, s.Settled(), "%s is failed but not settled", s)
		}
	}
}

func TestRetryAfter(t *testing.T) {
	now := time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		in   string
		want time.Duration
	}{
		{"absent", "", 0},
		{"seconds", "3", 3 * time.Second},
		{"padded seconds", " 2 ", 2 * time.Second},
		{"negative", "-1", 0},
		{"http date", now.Add(5 * time.Second).Format(http.TimeFormat), 5 * time.Second},
		{"past date", now.Add(-time.Minute).Format(http.TimeFormat), 0},
		{"garbage", "soon", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, retryAfter(tt.in, now))
		})
	}
}

func response(status int, header ...string) *http.Response {
	h := http.Header{}
	for i := 0; i+1 < len(header); i += 2 {
		h.Set(header[i], header[i+1])
	}
	return &http.Response{StatusCode: status, Header: h}
}

func TestRetryable(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		retry bool
		wait  time.Duration
	}{
		{"429 with Retry-After", NewStatusError(response(429, "Retry-After", "4"), errors.New("slow down")), true, 4 * time.Second},
		{"429 without Retry-After", NewStatusError(response(429), errors.New("slow down")), true, 0},
		{"500", NewStatusError(response(500), errors.New("boom")), true, 0},
		{"503", NewStatusError(response(503), errors.New("unavailable")), true, 0},
		{"404", NewStatusError(response(404), utils.ErrorNotFound404), false, 0},
		{"400", NewStatusError(response(400), errors.New("bad")), false, 0},
		{"no response", utils.ErrorInternalServerError, true, 0},
		{"timeout", fmt.Errorf("get: %w", utils.ErrorTimeoutAPICall), true, 0},
		{"plain", errors.New("decode"), false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retry, wait := Retryable(tt.err)
			assert.Equal(t, tt.retry, retry)
			assert.Equal(t, tt.wait, wait)
		})
	}
}

func TestStatusErrorKeepsTheShownError(t *testing.T) {
	err := NewStatusError(response(404), utils.ErrorNotFound404)
	assert.ErrorIs(t, err, utils.ErrorNotFound404)
	assert.Equal(t, utils.ErrorNotFound404.Error(), err.Error())

	shown := errors.New("no response")
	assert.Same(t, shown, NewStatusError(nil, shown))
}

func TestBuildError(t *testing.T) {
	err := &BuildError{
		Kind: "application", Name: "site", ID: 42, VersionID: "01JCA4M8",
		State: Error, TraceID: "abc", Err: ErrBuildFailed,
	}
	assert.Equal(t,
		`application "site" (id 42, version 01JCA4M8) ended in error: build did not produce a usable version (trace id abc)`,
		err.Error())
	assert.ErrorIs(t, errors.Join(err, errors.New("other")), ErrBuildFailed)

	var be *BuildError
	assert.True(t, errors.As(errors.Join(errors.New("other"), err), &be))
	assert.Equal(t, int64(42), be.ID)
}
