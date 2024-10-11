package schedule

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
)

func TestNewSchedule(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type args struct {
		fact *Factory
		name string
		kind string
	}
	tests := []struct {
		name    string
		args    args
		factory Factory
		wantErr bool
	}{
		{
			name: "happy road",
			args: args{
				name: "arthur-morgan",
				kind: DELETE_BUCKET,
			},
			factory: Factory{
				Dir:  func() config.DirPath { return config.DirPath{} },
				Join: func(elem ...string) string { return "" },
				Stat: os.Stat,
				IsNotExist: func(err error) bool {
					return true
				},
				WriteFile: func(name string, data []byte, perm os.FileMode) error {
					return nil
				},
				Unmarshal:     json.Unmarshal,
				MarshalIndent: json.MarshalIndent,
				Marshal:       json.Marshal,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewSchedule(nil, tt.args.name, tt.args.kind); (err != nil) != tt.wantErr {
				t.Errorf("NewSchedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecSchedules(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type args struct {
		factory *cmdutil.Factory
	}
	tests := []struct {
		name      string
		requests  []httpmock.Matcher
		responses []httpmock.Responder
		factory   Factory
		args      args
	}{
		{
			name: "happy road",
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodDelete, "v4/storage/buckets/arthur-morgan"),
			},
			responses: []httpmock.Responder{
				func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNoContent,
					}, nil
				},
			},
			factory: Factory{
				Dir:  func() config.DirPath { return config.DirPath{} },
				Join: func(elem ...string) string { return "" },
				Stat: os.Stat,
				IsNotExist: func(err error) bool {
					return false
				},
				WriteFile: func(name string, data []byte, perm os.FileMode) error {
					return nil
				},
				ReadFile: func(name string) ([]byte, error) {
					return []byte(`[
					{
					"name": "spirited-jaime",
					"time": "2024-10-11T14:37:34.908776421-03:00",
					"kind": "DeleteBucket"
					},
					{
					"name": "spirited-jaime2",
					"time": "2024-09-11T14:37:34.908776421-03:00",
					"kind": "DeleteBucket"
					}
					]`), nil
				},
				Unmarshal:     json.Unmarshal,
				MarshalIndent: json.MarshalIndent,
			},
		},
		{
			name: "second happy road",
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodDelete, "v4/storage/buckets/arthur-morgan"),
			},
			responses: []httpmock.Responder{
				func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNoContent,
					}, nil
				},
			},
			factory: Factory{
				Dir:  func() config.DirPath { return config.DirPath{} },
				Join: func(elem ...string) string { return "" },
				Stat: os.Stat,
				IsNotExist: func(err error) bool {
					return true
				},
				WriteFile: func(name string, data []byte, perm os.FileMode) error {
					return nil
				},
				ReadFile: func(name string) ([]byte, error) {
					return []byte(`[
					{
					"name": "spirited-jaime",
					"time": "2024-10-11T14:37:34.908776421-03:00",
					"kind": "DeleteBucket"
					},
					{
					"name": "spirited-jaime2",
					"time": "2024-09-11T14:37:34.908776421-03:00",
					"kind": "DeleteBucket"
					}
					]`), nil
				},
				Unmarshal:     json.Unmarshal,
				MarshalIndent: json.MarshalIndent,
				Marshal:       json.Marshal,
			},
		},
		{
			name: "error flow serializer func marshal for when path does not exist file in path ",
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodDelete, "v4/storage/buckets/arthur-morgan"),
			},
			responses: []httpmock.Responder{
				func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNoContent,
					}, nil
				},
			},
			factory: Factory{
				Dir:  func() config.DirPath { return config.DirPath{} },
				Join: func(elem ...string) string { return "" },
				Stat: os.Stat,
				IsNotExist: func(err error) bool {
					return true
				},
				WriteFile: func(name string, data []byte, perm os.FileMode) error {
					return nil
				},
				ReadFile: func(name string) ([]byte, error) {
					return []byte(`[
					{
					"name": "spirited-jaime",
					"time": "2024-10-11T14:37:34.908776421-03:00",
					"kind": "DeleteBucket"
					},
					{
					"name": "spirited-jaime2",
					"time": "2024-09-11T14:37:34.908776421-03:00",
					"kind": "DeleteBucket"
					}
					]`), nil
				},
				Unmarshal:     json.Unmarshal,
				MarshalIndent: json.MarshalIndent,
				Marshal: func(v any) ([]byte, error) {
					return []byte(""), errors.New("error Marshal")
				},
			},
		},
		{
			name: "error flow Write file for when path does not exist file in path",
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodDelete, "v4/storage/buckets/arthur-morgan"),
			},
			responses: []httpmock.Responder{
				func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNoContent,
					}, nil
				},
			},
			factory: Factory{
				Dir:  func() config.DirPath { return config.DirPath{} },
				Join: func(elem ...string) string { return "" },
				Stat: os.Stat,
				IsNotExist: func(err error) bool {
					return true
				},
				WriteFile: func(name string, data []byte, perm os.FileMode) error {
					return errors.New("error WriteFile")
				},
				ReadFile: func(name string) ([]byte, error) {
					return []byte(`[
					{
					"name": "spirited-jaime",
					"time": "2024-10-11T14:37:34.908776421-03:00",
					"kind": "DeleteBucket"
					},
					{
					"name": "spirited-jaime2",
					"time": "2024-09-11T14:37:34.908776421-03:00",
					"kind": "DeleteBucket"
					}
					]`), nil
				},
				Unmarshal:     json.Unmarshal,
				MarshalIndent: json.MarshalIndent,
				Marshal:       json.Marshal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InjectFactory(&tt.factory)

			mock := &httpmock.Registry{}
			for i, req := range tt.requests {
				mock.Register(req, tt.responses[i])
			}

			f, _, _ := testutils.NewFactory(mock)
			ExecSchedules(f)
		})
	}
}

func TestCheckIf24HoursPassed(t *testing.T) {
	tests := []struct {
		name   string
		passed time.Time
		want   bool
	}{
		{
			name:   "Scenario where 24 hours have already passed",
			passed: time.Now().Add(-25 * time.Hour),
			want:   true,
		},
		{
			name:   "Scenario where less than 24 hours have passed",
			passed: time.Now().Add(-23 * time.Hour),
			want:   false,
		},
		{
			name:   "Scenario where exactly 24 hours have passed",
			passed: time.Now().Add(-24 * time.Hour),
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckIf24HoursPassed(tt.passed); got != tt.want {
				t.Errorf("CheckIf24HoursPassed() = %v, want %v", got, tt.want)
			}
		})
	}
}
