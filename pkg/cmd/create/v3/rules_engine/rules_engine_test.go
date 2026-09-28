package rules_engine

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/create/rules_engine"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name      string
		args      []string
		request   httpmock.Matcher
		response  httpmock.Responder
		wantOut   string
		wantError string
		err       bool
	}{
		{
			name:     "success phase request",
			args:     []string{"--application-id", "1679423488", "--phase", "request", "--file", "./fixtures/create.json"},
			request:  httpmock.REST(http.MethodPost, "edge_applications/1679423488/rules_engine/request/rules"),
			response: httpmock.JSONFromFile("./fixtures/resp_phase_request.json"),
			wantOut:  fmt.Sprintf(msg.OutputSuccess, 210543),
			err:      false,
		},
		{
			name:     "success phase response",
			args:     []string{"--application-id", "1679423488", "--phase", "response", "--file", "./fixtures/create.json"},
			request:  httpmock.REST(http.MethodPost, "edge_applications/1679423488/rules_engine/response/rules"),
			response: httpmock.JSONFromFile("./fixtures/resp_phase_response.json"),
			wantOut:  fmt.Sprintf(msg.OutputSuccess, 210544),
			err:      false,
		},
		{
			name:      "error name empty",
			args:      []string{"--application-id", "1679423488", "--phase", "response", "--file", "./fixtures/create_name_empty.json"},
			request:   httpmock.REST(http.MethodPost, "edge_applications/1679423488/rules_engine/response/rules"),
			response:  httpmock.JSONFromFile("./fixtures/resp_phase_response.json"),
			wantError: msg.ErrorNameEmpty.Error(),
			err:       true,
		},
		{
			name:      "error conditional empty",
			args:      []string{"--application-id", "1679423488", "--phase", "response", "--file", "./fixtures/create_conditional_empty.json"},
			request:   httpmock.REST(http.MethodPost, "edge_applications/1679423488/rules_engine/response/rules"),
			response:  httpmock.JSONFromFile("./fixtures/resp_phase_response.json"),
			wantError: msg.ErrorConditionalEmpty.Error(),
			err:       true,
		},
		{
			name:      "error unmarshal file not exist",
			args:      []string{"--application-id", "1679423488", "--phase", "response", "--file", "./fixtures/no_exist.json"},
			request:   httpmock.REST(http.MethodPost, "edge_applications/1679423488/rules_engine/response/rules"),
			response:  httpmock.JSONFromFile("./fixtures/resp_phase_response.json"),
			wantError: utils.ErrorUnmarshalReader.Error(),
			err:       true,
		},
		{
			name:      "error api request create rules",
			args:      []string{"--application-id", "1679423488", "--phase", "request", "--file", "./fixtures/create.json"},
			request:   httpmock.REST(http.MethodPost, "edge_applications/1679423488/rules_engine/request/rules"),
			response:  httpmock.StatusStringResponse(http.StatusInternalServerError, "invalid"),
			wantError: "Failed to create the rule in Rules Engine: The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support. Check your settings and try again. If the error persists, contact Azion support.",
			err:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)
			f, outGot, _ := testutils.NewFactory(mock)

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if !tt.err && err == nil {
				require.Equal(t, tt.wantOut, outGot.String())
			} else {
				if err.Error() == tt.wantError {
					return
				}
				t.Fatal("Error: ", err)
			}
		})
	}
}

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name    string
		request sdk.CreateRulesEngineRequest
		wantErr bool
	}{
		{
			name:    "error name empty",
			request: sdk.CreateRulesEngineRequest{},
			wantErr: true,
		},
		{
			name: "error criteria null",
			request: sdk.CreateRulesEngineRequest{
				Name:     "no_empty",
				Criteria: nil,
			},
			wantErr: true,
		},
		{
			name: "error struct criteria conditional empty",
			request: sdk.CreateRulesEngineRequest{
				Name: "no_empty",
				Criteria: [][]sdk.RulesEngineCriteria{
					{{Conditional: ""}},
				},
			},
			wantErr: true,
		},
		{
			name: "error struct criteria variable empty",
			request: sdk.CreateRulesEngineRequest{
				Name: "no_empty",
				Criteria: [][]sdk.RulesEngineCriteria{
					{{Conditional: "no_empty", Variable: ""}},
				},
			},
			wantErr: true,
		},
		{
			name: "error struct criteria operator empty",
			request: sdk.CreateRulesEngineRequest{
				Name: "no_empty",
				Criteria: [][]sdk.RulesEngineCriteria{
					{{Conditional: "no_empty", Variable: "no_empty", Operator: ""}},
				},
			},
			wantErr: true,
		},
		{
			name: "error struct criteria variable empty",
			request: sdk.CreateRulesEngineRequest{
				Name: "no_empty",
				Criteria: [][]sdk.RulesEngineCriteria{
					{{Conditional: "no_empty", Variable: "no_empty", Operator: "no_empty"}},
				},
			},
			wantErr: true,
		},
		{
			name: "error behaviors null",
			request: sdk.CreateRulesEngineRequest{
				Name: "no_empty",
				Criteria: [][]sdk.RulesEngineCriteria{
					{{Conditional: "no_empty", Variable: "no_empty", Operator: "no_empty", InputValue: utils.PointerString("")}},
				},
				Behaviors: nil,
			},
			wantErr: true,
		},
		{
			name: "error struct string from behaviors field name empty",
			request: sdk.CreateRulesEngineRequest{
				Name: "no_empty",
				Criteria: [][]sdk.RulesEngineCriteria{
					{{Conditional: "no_empty", Variable: "no_empty", Operator: "no_empty", InputValue: utils.PointerString("")}},
				},
				Behaviors: []sdk.RulesEngineBehaviorEntry{
					sdk.RulesEngineBehaviorEntry{
						RulesEngineBehaviorString: &sdk.RulesEngineBehaviorString{
							Name: "",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error struct Object from behaviors field name empty",
			request: sdk.CreateRulesEngineRequest{
				Name: "no_empty",
				Criteria: [][]sdk.RulesEngineCriteria{
					{{Conditional: "no_empty", Variable: "no_empty", Operator: "no_empty", InputValue: utils.PointerString("")}},
				},
				Behaviors: []sdk.RulesEngineBehaviorEntry{
					sdk.RulesEngineBehaviorEntry{
						RulesEngineBehaviorObject: &sdk.RulesEngineBehaviorObject{},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRequest(tt.request); (err != nil) != tt.wantErr {
				t.Errorf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
