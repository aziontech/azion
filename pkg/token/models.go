package token

import (
	"io"
	"net/http"
)

const settingsFilename = "settings.toml"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Token struct {
	endpoint string
	client   HTTPClient
	filePath string
	valid    bool
	out      io.Writer
}

type Settings struct {
	Token string
	UUID  string
}

type Config struct {
	Client HTTPClient
	Out    io.Writer
}

type Response struct {
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}
