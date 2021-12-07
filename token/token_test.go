package token

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const validToken = "rightToken"

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

func Test_Validate(t *testing.T) {
	t.Run("invalid token", func(t *testing.T) {
		client := &MockClient{}

		token := NewToken(client)

		tokenString := "thisIsNotTheValidToken"
		valid, _ := token.Validate(&tokenString)

		if valid {
			t.Errorf("token is valid; expected to be invalid")
		}
	})

	t.Run("valid token", func(t *testing.T) {
		client := &MockClient{}

		token := NewToken(client)

		tokenString := "rightToken"
		valid, _ := token.Validate(&tokenString)

		if !valid {
			t.Errorf("token is not valid; expected to be valid")
		}
	})

}

func Test_Save(t *testing.T) {
	t.Run("save token to disk", func(t *testing.T) {
		client := &MockClient{}
		token := NewToken(client)
		token.token = "TeSt"
		if token.Save() != nil {
			t.Errorf("token saved " + token.token)
		}

	})
}

func Test_ReadFromDisk(t *testing.T) {
	t.Run("read token from disk", func(t *testing.T) {
		client := &MockClient{}
		token := NewToken(client)
		token.token = "TeSt"
		dToken, _ := token.ReadFromDisk()
		if dToken != token.token {
			t.Errorf("token from disk differs from test: " + token.token)
		}

	})
}
