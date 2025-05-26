package edge_storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/zRedShift/mimemagic"
	"go.uber.org/zap/zapcore"
)

func TestNewObjects(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name     string
		fact     factoryObjects
		request  httpmock.Matcher
		response httpmock.Responder
		args     []string
		output   string
		wantErr  bool
		Err      string
	}{
		{
			name: "create new object command object of the edge-storage",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              utils.AskInput,
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--object-key", "nomedoobject", "--source", "fixtures/index.html"},
			output:   msg.OUTPUT_CREATE_OBJECT,
		},
		{
			name: "create new object command object of the edge-storage using flag --file",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              utils.AskInput,
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--file", "fixtures/create_object.json"},
			output:   msg.OUTPUT_CREATE_OBJECT,
		},
		{
			name: "input file json err --file",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              utils.AskInput,
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "/"),
			response: httpmock.JSONFromFile("/"),
			args:     []string{"--file", "fixtures/create_error.json"},
			Err:      utils.ErrorUnmarshalReader.Error(),
		},
		{
			name: "failed internal error status 500",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              utils.AskInput,
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request: httpmock.REST(http.MethodPost, "v4/storage/buckets"),
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, utils.ErrorInternalServerError
			},
			args: []string{"--file", "fixtures/create_object.json"},
			Err:  fmt.Sprintf(msg.ERROR_CREATE_OBJECT, "The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support"),
		},
		{
			name: "error ask input bucket name",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              func(msg string) (string, error) { return "", utils.ErrorParseResponse },
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--object-key", "nomedoobject", "--source", "fixtures/index.html"},
			Err:      utils.ErrorParseResponse.Error(),
			wantErr:  true,
		},
		{
			name: "success ask input bucket name",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              func(msg string) (string, error) { return "nomedobucket", nil },
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--object-key", "nomedoobject", "--source", "fixtures/index.html"},
			output:   msg.OUTPUT_CREATE_OBJECT,
		},
		{
			name: "error ask input objec key",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              func(msg string) (string, error) { return "", utils.ErrorParseResponse },
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--source", "fixtures/index.html"},
			Err:      utils.ErrorParseResponse.Error(),
			wantErr:  true,
		},
		{
			name: "success ask input object key",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              func(msg string) (string, error) { return "nomedoobject", nil },
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--source", "fixtures/index.html"},
			output:   msg.OUTPUT_CREATE_OBJECT,
		},
		{
			name: "error ask input source",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              func(msg string) (string, error) { return "", utils.ErrorParseResponse },
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--object-key", "nomedoobject"},
			Err:      utils.ErrorParseResponse.Error(),
			wantErr:  true,
		},
		{
			name: "success ask input source",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              func(msg string) (string, error) { return "fixtures/index.html", nil },
				Open:                  os.Open,
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--object-key", "nomedoobject"},
			output:   msg.OUTPUT_CREATE_OBJECT,
		},
		{
			name: "error func internal Getwd",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              utils.AskInput,
				Open:                  os.Open,
				Getwd:                 func() (dir string, err error) { return "", utils.ErrorInternalServerError },
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--object-key", "nomedoobject", "--source", "fixtures/index.html"},
			Err:      utils.ErrorInternalServerError.Error(),
			wantErr:  true,
		},
		{
			name: "error open file",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath:         mimemagic.MatchFilePath,
				AskInput:              utils.AskInput,
				Open:                  func(name string) (*os.File, error) { return nil, errors.New("error open file") },
				Getwd:                 os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--object-key", "nomedoobject", "--source", "fixtures/index.html"},
			Err:      "error open file",
			wantErr:  true,
		},
		{
			name: "error open file with mimetype",
			fact: factoryObjects{
				FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
				Join:                  filepath.Join,
				MatchFilePath: func(path string, limAndPref ...int) (mimemagic.MediaType, error) {
					return mimemagic.MediaType{}, errors.New("error open file")
				},
				AskInput: utils.AskInput,
				Open:     os.Open,
				Getwd:    os.Getwd,
			},
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--object-key", "nomedoobject", "--source", "fixtures/index.html"},
			Err:      "error open file",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)
			f, out, _ := testutils.NewFactory(mock)
			tt.fact.Factory = f
			cmd := commandObjects(&tt.fact)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				if !strings.EqualFold(tt.Err, err.Error()) {
					t.Errorf("Error expected: %s got: %s", tt.Err, err.Error())
				}
			} else {
				assert.Equal(t, tt.output, out.String())
			}
		})
	}
}

func TestNewFactoryObjects(t *testing.T) {
	mock := &httpmock.Registry{}
	f, _, _ := testutils.NewFactory(mock)
	NewFactoryObjects(f)
}
