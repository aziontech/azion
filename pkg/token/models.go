package token

import (
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Token struct {
	Endpoint string
	client   HTTPClient
	filePath string
	valid    bool
	out      io.Writer
}

type UserInfo struct {
	Results struct {
		ClientID string `json:"client_id"`
		Email    string `json:"email"`
	} `json:"results"`
}

type Profile struct {
	Name string
}

type Settings struct {
	Token                      string
	UUID                       string
	LastCheck                  time.Time
	LastVulcanVersion          string
	AuthorizeMetricsCollection int
	ClientId                   string
	Email                      string
	ContinuationToken          string
	S3AccessKey                string
	S3SecretKey                string
	S3Bucket                   string
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

// S3Credentials represents the S3 credentials for a specific bucket
type S3Credentials struct {
	S3AccessKey string `toml:"s3_access_key"`
	S3SecretKey string `toml:"s3_secret_key"`
}

// CredentialsMap is a map of bucket names to their S3 credentials
type CredentialsMap map[string]S3Credentials
