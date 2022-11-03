package update

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
)

var responseBody = `
{
    "id": 1986,
    "name": "{name}",
    "updated_at": "2021-02-01T10:00:00Z",
    "last_editor": "pepe@azion.com",
    "active": false,
    "bound_nodes": 0,
    "permissions": [
        "read",
        "write"
    ]
}
`

func TestUpdate(t *testing.T) {
	t.Run("update service with name", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_services/1234"),
			func(req *http.Request) (*http.Response, error) {
				request := &sdk.UpdateServiceRequest{}
				body, _ := io.ReadAll(req.Body)
				_ = json.Unmarshal(body, request)
				response := strings.ReplaceAll(responseBody, "{name}", *request.Name)

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
		cmd.SetArgs([]string{"--service-id", "1234", "--name", "thunderstruck"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
