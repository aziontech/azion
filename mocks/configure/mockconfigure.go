package configure

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

const validToken = "d813e2a28bda82f46e6403afdb8ed8a2516df1df"

type MockClient struct {
}

func (m MockClient) Get(req *http.Request) (*http.Response, error) {
	token := req.URL.Query().Get("token")
	valid := validToken == token
	body := `{"valid": ` + strconv.FormatBool(valid) + `}`

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil

}
