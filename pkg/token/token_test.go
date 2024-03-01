package token

import (
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func Test_Validate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

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
		valid, _, _ := token.Validate(&tokenString)

		if valid {
			t.Errorf("Validate() = %v; want false", valid)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "token/user/me"),
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
		valid, _, _ := token.Validate(&tokenString)

		if !valid {
			t.Errorf("Validate() = %v; want true", valid)
		}
	})

}

func Test_Save(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	errSetPath := config.SetPath("/tmp/testazion/test.toml")
	if errSetPath != nil {
		t.Fatalf("SetPath() error: %s;", errSetPath.Error())
	}

	t.Run("save token to disk", func(t *testing.T) {
		token, err := New(&Config{
			Out: os.Stdout,
		})

		if err != nil {
			t.Fatalf("New() = %v; want nil", err)
		}

		settings := Settings{
			Token: "asdfdsafsdfasd",
			UUID:  "asdfadsfads",
		}

		b, err := toml.Marshal(settings)
		require.NoError(t, err)

		if _, err := token.Save(b); err != nil {
			t.Fatalf("Save() = %v; want nil", err)
		}
	})
}
