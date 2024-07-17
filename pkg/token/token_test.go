package token

import (
	"net/http"
	"os"
	"path/filepath"
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

func Test_Create(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("create token", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "token/tokens"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		token, err := New(&Config{
			Client: &http.Client{Transport: mock},
			Out:    os.Stdout,
		})
		if err != nil {
			t.Fatalf("NewToken() = %v; want nil", err)
		}

		token.endpoint = "http://api.azion.net/token"
		response, err := token.Create("base64Credentials")
		if err != nil {
			t.Fatalf("Create() = %v; want nil", err)
		}

		if response.Token != "123321" {
			t.Errorf("Create() = %v; want 123321", response.Token)
		}
	})
}

func Test_WriteSettings(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("write settings to disk", func(t *testing.T) {
		settings := Settings{
			Token: "tokenValue",
			UUID:  "uuidValue",
		}

		errSetPath := config.SetPath("/tmp/testazion/test.toml")
		if errSetPath != nil {
			t.Fatalf("SetPath() error: %s;", errSetPath.Error())
		}

		err := WriteSettings(settings)
		if err != nil {
			t.Fatalf("WriteSettings() = %v; want nil", err)
		}

		dir, err := config.Dir()
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(dir.Dir, dir.Settings))
		require.NoError(t, err)

		var readSettings Settings
		err = toml.Unmarshal(data, &readSettings)
		require.NoError(t, err)

		if readSettings.Token != settings.Token || readSettings.UUID != settings.UUID {
			t.Errorf("WriteSettings() wrote %v; want %v", readSettings, settings)
		}
	})
}

func Test_ReadSettings(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("read settings from disk", func(t *testing.T) {
		errSetPath := config.SetPath("/tmp/testazion/test.toml")
		if errSetPath != nil {
			t.Fatalf("SetPath() error: %s;", errSetPath.Error())
		}

		expectedSettings := Settings{
			Token: "tokenValue",
			UUID:  "uuidValue",
		}

		err := WriteSettings(expectedSettings)
		if err != nil {
			t.Fatalf("WriteSettings() error: %s", err)
		}

		settings, err := ReadSettings()
		if err != nil {
			t.Fatalf("ReadSettings() = %v; want nil", err)
		}

		if settings.Token != expectedSettings.Token || settings.UUID != expectedSettings.UUID {
			t.Errorf("ReadSettings() read %v; want %v", settings, expectedSettings)
		}
	})

	t.Run("read settings from non-existing file", func(t *testing.T) {
		errSetPath := config.SetPath("/tmp/testazion/nonexistent.toml")
		if errSetPath != nil {
			t.Fatalf("SetPath() error: %s;", errSetPath.Error())
		}

		settings, err := ReadSettings()
		if err == nil {
			t.Fatalf("ReadSettings() error = nil; want non-nil error")
		}

		if settings != (Settings{}) {
			t.Errorf("ReadSettings() = %v; want empty settings", settings)
		}
	})
}
