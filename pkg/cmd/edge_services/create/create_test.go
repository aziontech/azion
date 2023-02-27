package create

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/stretchr/testify/require"
	msg "github.com/aziontech/azion-cli/messages/edge_services"
)

var responseBody = `
{
    "id": 1753,
    "name": "{name}",
    "updated_at": "2021-12-16T01:10:07Z",
    "last_editor": "crazy.ape@azion.com",
    "active": false,
    "bound_nodes": 0,
    "permissions": [
        "read",
        "write"
    ]
}
`

func TestCreate(t *testing.T) {
	t.Run("invalid service", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/"),
			httpmock.StatusStringResponse(http.StatusUnprocessableEntity, "Invalid name"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", ""})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("without passing name", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, msg.ErrorMandatoryName)
	})

	t.Run("create service with name", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/"),
			func(req *http.Request) (*http.Response, error) {
				request := &sdk.CreateServiceRequest{}
				body, _ := io.ReadAll(req.Body)
				_ = json.Unmarshal(body, request)

				response := strings.ReplaceAll(responseBody, "{name}", request.Name)

				return &http.Response{StatusCode: http.StatusCreated,
					Request: req,
					Body:    io.NopCloser(strings.NewReader(response)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "BIRL"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
