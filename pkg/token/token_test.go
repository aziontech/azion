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

		token := New(&Config{
			Client: &http.Client{Transport: mock},
			Out:    os.Stdout,
		})

		token.Endpoint = "http://api.azion.net/token"
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

		token := New(&Config{
			Client: &http.Client{Transport: mock},
			Out:    os.Stdout,
		})

		token.Endpoint = "http://api.azion.net/token"
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
		token := New(&Config{
			Out: os.Stdout,
		})

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

		token := New(&Config{
			Client: &http.Client{Transport: mock},
			Out:    os.Stdout,
		})

		token.Endpoint = "http://api.azion.net/token"
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

		err := WriteSettings(settings, "test")
		if err != nil {
			t.Fatalf("WriteSettings() = %v; want nil", err)
		}

		dir := config.Dir()
		data, err := os.ReadFile(filepath.Join(dir.Dir, "test", dir.Settings))
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

		err := WriteSettings(expectedSettings, "test")
		if err != nil {
			t.Fatalf("WriteSettings() error: %s", err)
		}

		settings, err := ReadSettings("test")
		if err != nil {
			t.Fatalf("ReadSettings() = %v; want nil", err)
		}

		if settings.Token != expectedSettings.Token || settings.UUID != expectedSettings.UUID {
			t.Errorf("ReadSettings() read %v; want %v", settings, expectedSettings)
		}
	})

	t.Run("read settings from non-existing file creates default", func(t *testing.T) {
		errSetPath := config.SetPath("/tmp/testazion/nonexistent.toml")
		if errSetPath != nil {
			t.Fatalf("SetPath() error: %s;", errSetPath.Error())
		}

		settings, err := ReadSettings("nonexistent")
		if err != nil {
			t.Fatalf("ReadSettings() error = %v; want nil", err)
		}

		// Should return default empty settings
		expectedSettings := Settings{}
		if settings != expectedSettings {
			t.Errorf("ReadSettings() = %v; want %v", settings, expectedSettings)
		}
	})
}

func Test_ReadWriteCredentials(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("write and read credentials", func(t *testing.T) {
		errSetPath := config.SetPath("/tmp/testazioncreds/test.toml")
		if errSetPath != nil {
			t.Fatalf("SetPath() error: %s;", errSetPath.Error())
		}

		bucketName := "test-bucket-123"
		expectedCreds := S3Credentials{
			S3AccessKey: "test-access-key",
			S3SecretKey: "test-secret-key",
		}

		// Write credentials
		err := SaveCredentialsForBucket("test", bucketName, expectedCreds)
		if err != nil {
			t.Fatalf("SaveCredentialsForBucket() error: %s", err)
		}

		// Read credentials
		creds, exists, err := GetCredentialsForBucket("test", bucketName)
		if err != nil {
			t.Fatalf("GetCredentialsForBucket() error: %s", err)
		}

		if !exists {
			t.Fatalf("GetCredentialsForBucket() exists = %v; want true", exists)
		}

		if creds.S3AccessKey != expectedCreds.S3AccessKey || creds.S3SecretKey != expectedCreds.S3SecretKey {
			t.Errorf("GetCredentialsForBucket() = %v; want %v", creds, expectedCreds)
		}
	})

	t.Run("read credentials from non-existing bucket", func(t *testing.T) {
		errSetPath := config.SetPath("/tmp/testazioncreds/nonexistent.toml")
		if errSetPath != nil {
			t.Fatalf("SetPath() error: %s;", errSetPath.Error())
		}

		creds, exists, err := GetCredentialsForBucket("nonexistent", "nonexistent-bucket")
		if err != nil {
			t.Fatalf("GetCredentialsForBucket() error = %v; want nil", err)
		}

		if exists {
			t.Fatalf("GetCredentialsForBucket() exists = %v; want false", exists)
		}

		if creds.S3AccessKey != "" || creds.S3SecretKey != "" {
			t.Errorf("GetCredentialsForBucket() = %v; want empty credentials", creds)
		}
	})

	t.Run("read credentials from non-existing file returns empty map", func(t *testing.T) {
		errSetPath := config.SetPath("/tmp/testazioncreds/newprofile/test.toml")
		if errSetPath != nil {
			t.Fatalf("SetPath() error: %s;", errSetPath.Error())
		}

		credentials, err := ReadCredentials("newprofile")
		if err != nil {
			t.Fatalf("ReadCredentials() error = %v; want nil", err)
		}

		if credentials == nil {
			t.Fatalf("ReadCredentials() = nil; want empty map")
		}

		if len(credentials) != 0 {
			t.Errorf("ReadCredentials() = %v; want empty map", credentials)
		}
	})

	t.Run("write multiple buckets and read them", func(t *testing.T) {
		errSetPath := config.SetPath("/tmp/testazioncreds/multi/test.toml")
		if errSetPath != nil {
			t.Fatalf("SetPath() error: %s;", errSetPath.Error())
		}

		bucket1 := "bucket-one"
		creds1 := S3Credentials{
			S3AccessKey: "access-key-1",
			S3SecretKey: "secret-key-1",
		}

		bucket2 := "bucket-two"
		creds2 := S3Credentials{
			S3AccessKey: "access-key-2",
			S3SecretKey: "secret-key-2",
		}

		// Write both credentials
		err := SaveCredentialsForBucket("multi", bucket1, creds1)
		if err != nil {
			t.Fatalf("SaveCredentialsForBucket() error for bucket1: %s", err)
		}

		err = SaveCredentialsForBucket("multi", bucket2, creds2)
		if err != nil {
			t.Fatalf("SaveCredentialsForBucket() error for bucket2: %s", err)
		}

		// Verify both credentials exist and are correct
		readCreds1, exists1, err := GetCredentialsForBucket("multi", bucket1)
		if err != nil {
			t.Fatalf("GetCredentialsForBucket() error for bucket1: %s", err)
		}
		if !exists1 {
			t.Fatalf("GetCredentialsForBucket() exists for bucket1 = %v; want true", exists1)
		}
		if readCreds1.S3AccessKey != creds1.S3AccessKey || readCreds1.S3SecretKey != creds1.S3SecretKey {
			t.Errorf("GetCredentialsForBucket() bucket1 = %v; want %v", readCreds1, creds1)
		}

		readCreds2, exists2, err := GetCredentialsForBucket("multi", bucket2)
		if err != nil {
			t.Fatalf("GetCredentialsForBucket() error for bucket2: %s", err)
		}
		if !exists2 {
			t.Fatalf("GetCredentialsForBucket() exists for bucket2 = %v; want true", exists2)
		}
		if readCreds2.S3AccessKey != creds2.S3AccessKey || readCreds2.S3SecretKey != creds2.S3SecretKey {
			t.Errorf("GetCredentialsForBucket() bucket2 = %v; want %v", readCreds2, creds2)
		}
	})
}
