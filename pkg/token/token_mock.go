package token

type TokenMock struct {
}

func (m *TokenMock) Validate(token *string) (bool, UserInfo, error) {
	return true, UserInfo{}, nil
}

func (m *TokenMock) Save(b []byte) (string, error) {
	return "", nil
}

func (m *TokenMock) Create(b64 string) (*Response, error) {
	return &Response{}, nil
}
