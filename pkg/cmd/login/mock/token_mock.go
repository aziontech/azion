package mock

import (
	"net/http"

	"github.com/aziontech/azion-cli/pkg/token"
)

type TokenMock struct{}

func Validate(t *string) (bool, token.UserInfo, error) {
	return false, token.UserInfo{}, nil
}

func Save(b []byte) (string, error) {
	return "", nil
}

func Create(b64 string) (*http.Response, error) {
	return &http.Response{}, nil
}
