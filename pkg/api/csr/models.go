package csr

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.CertificateSigningRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}
