package create

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	errmsg "github.com/aziontech/azion-cli/messages/edge_services"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var resposeBody = `
{
    "id": 82706,
    "name": "{name}",
    "type": "{type}",
    "content": "{content}",
    "content_type": "{content_type}"
}
`

func buildResponseContent(req *http.Request) string {
	request := &sdk.CreateResourceRequest{}
	body, _ := ioutil.ReadAll(req.Body)
	_ = json.Unmarshal(body, request)

	response := strings.ReplaceAll(resposeBody, "{name}", request.Name)
	response = strings.ReplaceAll(response, "{type}", request.Trigger)
	response = strings.ReplaceAll(response, "{content}", request.Content)
	response = strings.ReplaceAll(response, "{content_type}", request.ContentType)

	return response
}

func TestCreate(t *testing.T) {
	t.Run("create text resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/1234/resources"),
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusCreated,
					Request: req,
					Body:    ioutil.NopCloser(strings.NewReader(buildResponseContent(req))),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		f, stdout, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		_, _ = contentFile.Write([]byte("insert your text here"))

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--service-id", "1234", "--name", "/tmp/testando.txt", "--content-type", "text", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, `Created Resource with ID 82706
`, stdout.String())

	})

	t.Run("create script resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/1234/resources"),
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusCreated,
					Request:    req,
					Body:       ioutil.NopCloser(strings.NewReader(buildResponseContent(req))),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		f, _, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		_, _ = contentFile.Write([]byte("#!/bin/sh"))

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--service-id", "1234", "--name", "/tmp/bomb.sh", "--trigger", "Install", "--content-type", "shellscript", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("create script resource without trigger", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-s", "1234", "--name", "/tmp/bomb.sh", "--content-type", "shellscript", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, errmsg.ErrorInvalidResourceTrigger)
	})

	t.Run("create resource without content file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-s", "1234", "--name", "/tmp/a.txt", "--content-type", "text"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.EqualError(t, err, "You must provide --name, --content-type and --content-file flags when --in flag is not sent")
	})

	t.Run("content file cannot be empty", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/1234/resources"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Bad Request"),
		)
		f, _, _ := testutils.NewFactory(mock)

		contentFile, _ := os.CreateTemp("", "content.txt")

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--service-id", "1234", "--name", "/tmp/a.txt", "--content-type", "text", "--content-file", contentFile.Name()})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)
		_, err := cmd.ExecuteC()
		require.EqualError(t, err, "The file’s content is empty. Provide a path and file with valid content and try the command again. Use the flags -h or --help with a command or subcommand to display more information and try again")
	})
}
