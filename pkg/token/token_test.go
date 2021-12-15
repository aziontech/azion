package token

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

const validToken = "rightToken"

type MockClient struct {
}

func (m MockClient) Do(req *http.Request) (*http.Response, error) {
	header := req.Header.Get("Authorization")

	// Authorization: token <token>
	split := strings.Split(header, "token ")

	valid := validToken == split[1]

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

		token, err := NewToken(&Config{
			Client: client,
			Out:    os.Stdout,
		})
		if err != nil {
			t.Fatalf("NewToken() = %v; want nil", err)
		}

		tokenString := "thisIsNotTheValidToken"
		valid, _ := token.Validate(&tokenString)

		if valid {
			t.Errorf("Validate() = %v; want false", valid)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		client := &MockClient{}

		token, err := NewToken(&Config{
			Client: client,
			Out:    os.Stdout,
		})
		if err != nil {
			t.Fatalf("NewToken() = %v; want nil", err)
		}

		tokenString := "rightToken"
		valid, _ := token.Validate(&tokenString)

		if !valid {
			t.Errorf("Validate() = %v; want true", valid)
		}
	})

}

func Test_Save(t *testing.T) {
	t.Run("save token to disk", func(t *testing.T) {
		client := &MockClient{}

		token, err := NewToken(&Config{
			Client: client,
			Out:    os.Stdout,
		})
		if err != nil {
			t.Fatalf("NewToken() = %v; want nil", err)
		}

		token.filepath = "/tmp/azion/credentials"

		token.token = "TeSt"
		if err := token.Save(); err != nil {
			t.Fatalf("Save() = %v; want nil", err)
		}

	})
}

func Test_ReadFromDisk(t *testing.T) {
	t.Run("read token from disk", func(t *testing.T) {
		client := &MockClient{}
		token, err := NewToken(&Config{
			Client: client,
			Out:    os.Stdout,
		})
		if err != nil {
			t.Fatalf("NewToken() = %v; want nil", err)
		}

		token.filepath = "/tmp/azion/credentials"

		token.token = "TeSt"
		if err := token.Save(); err != nil {
			t.Fatalf("Save() = %v; want nil", err)
		}

		dToken, _ := token.ReadFromDisk()
		if dToken != token.token {
			t.Errorf("ReadFromDisk() = %v; want %v", dToken, token.token)
		}

	})
}
