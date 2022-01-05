package update

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/stretchr/testify/assert"
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
				body, _ := ioutil.ReadAll(req.Body)
				_ = json.Unmarshal(body, request)
				response := strings.ReplaceAll(responseBody, "{name}", *request.Name)

				return &http.Response{StatusCode: http.StatusCreated,
					Request: req,
					Body:    ioutil.NopCloser(strings.NewReader(response)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.PersistentFlags().BoolP("verbose", "v", false, "")
		cmd.SetArgs([]string{"1234", "--name", "thunderstruck"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("update service with name being verbose", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_services/1234"),
			func(req *http.Request) (*http.Response, error) {
				request := &sdk.UpdateServiceRequest{}
				body, _ := ioutil.ReadAll(req.Body)
				_ = json.Unmarshal(body, request)
				response := strings.ReplaceAll(responseBody, "{name}", *request.Name)

				return &http.Response{StatusCode: http.StatusCreated,
					Request: req,
					Body:    ioutil.NopCloser(strings.NewReader(response)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.PersistentFlags().BoolP("verbose", "v", false, "")
		cmd.SetArgs([]string{"1234", "--name", "thunderstruck", "-v"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, "ID: 1986\nName: thunderstruck\nUpdated at: 2021-02-01T10:00:00Z\nLast Editor: pepe@azion.com\nActive: false\nBound Nodes: 0\nPermissions: [read write]\n", stdout.String())

	})
}
