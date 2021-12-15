package token

import (
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
)

func Test_Validate(t *testing.T) {
	t.Run("invalid token", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "token"),
			httpmock.StatusStringResponse(http.StatusUnauthorized, "{}"),
		)

		token, err := New(&Config{
			Client: &http.Client{Transport: mock},
			Out:    os.Stdout,
		})
		if err != nil {
			t.Fatalf("NewToken() = %v; want nil", err)
		}

		token.endpoint = "http://api.azion.net/token"
		tokenString := "thisIsNotTheValidToken"
		valid, _ := token.Validate(&tokenString)

		if valid {
			t.Errorf("Validate() = %v; want false", valid)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "token"),
			httpmock.StatusStringResponse(http.StatusOK, "{}"),
		)

		token, err := New(&Config{
			Client: &http.Client{Transport: mock},
			Out:    os.Stdout,
		})
		if err != nil {
			t.Fatalf("NewToken() = %v; want nil", err)
		}

		token.endpoint = "http://api.azion.net/token"
		tokenString := "rightToken"
		valid, _ := token.Validate(&tokenString)

		if !valid {
			t.Errorf("Validate() = %v; want true", valid)
		}
	})

}

func Test_Save(t *testing.T) {
	t.Run("save token to disk", func(t *testing.T) {
		token := &Token{
			out:      os.Stdout,
			filepath: "/tmp/azion/credentials",
			token:    "TeST",
		}
		if err := token.Save(); err != nil {
			t.Fatalf("Save() = %v; want nil", err)
		}
	})
}

func Test_ReadFromDisk(t *testing.T) {
	t.Run("read token from disk", func(t *testing.T) {
		token := &Token{
			out:      os.Stdout,
			filepath: "/tmp/azion/credentials",
			token:    "TeST",
		}
		if err := token.Save(); err != nil {
			t.Fatalf("Save() = %v; want nil", err)
		}

		dToken, _ := token.ReadFromDisk()
		if dToken != token.token {
			t.Errorf("ReadFromDisk() = %v; want %v", dToken, token.token)
		}
	})
}
