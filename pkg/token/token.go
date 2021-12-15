package token

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Token struct {
	endpoint string
	client   HTTPClient
	filepath string
	token    string
	valid    bool
	out      io.Writer
}

var AuthEndpoint string

const credentialsFilename = "credentials"

type Config struct {
	Client HTTPClient
	Out    io.Writer
}

func New(c *Config) (*Token, error) {
	dir, err := TokenDir()
	if err != nil {
		return nil, err
	}

	return &Token{
		client:   c.Client,
		endpoint: AuthEndpoint,
		filepath: filepath.Join(dir, credentialsFilename),
		out:      c.Out,
	}, nil
}

func (t *Token) Validate(token *string) (bool, error) {
	req, err := http.NewRequest("GET", t.endpoint, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Accept", "application/json; version=3")
	req.Header.Add("Authorization", "token "+*token)

	resp, err := t.client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	t.token = *token
	t.valid = true

	return true, nil
}

func (t *Token) Save() error {
	fbyte := []byte(t.token)

	err := os.MkdirAll(filepath.Dir(t.filepath), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(t.filepath, fbyte, 0600)
	if err != nil {
		return err
	}

	fmt.Fprintln(t.out, "Token saved in "+t.filepath)

	return nil
}

func (t *Token) ReadFromDisk() (string, error) {
	filedata, err := os.ReadFile(t.filepath)
	if err != nil {
		return "", err
	}

	return string(filedata[:]), nil
}

func TokenDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".azion"), nil
}
