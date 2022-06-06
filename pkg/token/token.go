package token

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/pkg/constants"
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
		endpoint: constants.AuthURL,
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

	// TODO: provide a better description after token is saved
	fmt.Fprintf(t.out, "Token saved in %v\n", t.filepath)
	fmt.Fprintln(t.out, "This token will be used by default when calling any command")

	return nil
}

func ReadFromDisk() (string, error) {
	dir, err := TokenDir()
	if err != nil {
		return "", fmt.Errorf("failed to get token dir: %w", err)
	}

	filedata, err := os.ReadFile(filepath.Join(dir, credentialsFilename))
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
