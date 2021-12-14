package token

import (
	"fmt"
	"net/http"
	"os"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Token struct {
	client HTTPClient
	token  string
	valid  bool
}

var AUTH_ENDPOINT string

func NewToken(c HTTPClient) *Token {
	return &Token{c, "", false}
}

func (t *Token) Validate(token *string) (bool, error) {
	urlAuthentication := AUTH_ENDPOINT

	req, err := http.NewRequest("GET", urlAuthentication, nil)
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
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirname = dirname + "/.azion/"
	err = os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}

	filename := dirname + "credentials"
	err = os.WriteFile(filename, fbyte, 0600)
	if err != nil {
		return err
	}

	fmt.Println("Token saved in " + filename)
	return nil
}

func (t *Token) ReadFromDisk() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filename := dirname + "/.azion/credentials"
	filedata, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(filedata[:]), nil
}
