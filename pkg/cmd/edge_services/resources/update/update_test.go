package update

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/stretchr/testify/require"
)

var responseBody = `
{
	"id": 666,
	"name": "{name}",
	"type": "Install",
	"content": "Parangaricutirimírruaro",
	"content_type": "Shell Script"
  }
`

func TestUpdate(t *testing.T) {

	t.Run("not all arguments were sent", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_services/1234/resources/666"),
			httpmock.StringResponse("Error: You must provide a service_id and a resource_id as arguments. Use -h or --help for more information"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, errmsg.ErrorMissingArgumentUpdateResource)
	})

	t.Run("no flag was sent", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_services/1234/resources/666"),
			httpmock.StringResponse("Error: You must provide at least one value in update. Use -h or --help for more information"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "666"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, utils.ErrorUpdateNoFlagsSent)
	})

	t.Run("update resource with name", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_services/1234/resources/666"),
			func(req *http.Request) (*http.Response, error) {
				request := &sdk.UpdateResourceRequest{}
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
		cmd.SetArgs([]string{"1234", "666", "--name", "BIRL"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

	})

	t.Run("update resource with all felds", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_services/1234/resources/666"),
			func(req *http.Request) (*http.Response, error) {
				request := &sdk.UpdateResourceRequest{}
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

		contentFile, _ := os.CreateTemp("", "content.txt")

		_, _ = contentFile.Write([]byte("This content is made for testing purposes"))

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"1234", "666", "--name", "BIRL", "--trigger", "Install", "--content-type", "shellscript", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

}
