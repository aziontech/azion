package token

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

const validToken = "rightToken"

type MockClient struct {
}

func (m MockClient) Do(req *http.Request) (*http.Response, error) {

	header := req.Header.Get("Authorization")

	// Authorization: token <token>
	splitted := strings.Split(header, "token ")

	valid := validToken == splitted[1]

	var statusCode int
	if valid {
		statusCode = http.StatusOK
	} else {
		statusCode = http.StatusUnauthorized
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(nil),
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
