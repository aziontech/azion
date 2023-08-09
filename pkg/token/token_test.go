package token

import (
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/config"
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
	config.SetPath("/tmp/testazion")

	t.Run("save token to disk", func(t *testing.T) {
		token, err := New(&Config{
			Out: os.Stdout,
		})
		token.token = "TeST"

		if err != nil {
			t.Fatalf("New() = %v; want nil", err)
		}

		if err := token.Save(); err != nil {
			t.Fatalf("Save() = %v; want nil", err)
		}
	})
}
