package token

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func init() {
	os.Setenv("SSO_MODE", "development")
}

const validToken = "certo"

type MockClient struct {
}

func (m MockClient) Do(req *http.Request) (*http.Response, error) {
	token := req.URL.Query().Get("token")
	valid := validToken == token
	body := `{"valid": ` + strconv.FormatBool(valid) + `}`

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func Test_Authenticate(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		client := &MockClient{}

		token := NewToken(client)

		valid, _ := token.Validate("certo")

		if !valid {
			t.Fatalf("token is not valid; expected to be valid")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		client := &MockClient{}

		token := NewToken(client)

		valid, _ := token.Validate("elevador")

		if valid {
			t.Fatalf("token is valid; expected to be invalid")
		}
	})
}
