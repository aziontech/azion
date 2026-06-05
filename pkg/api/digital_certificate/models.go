package digitalcertificate

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.Certificate
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedCertificate
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}

type RequestRequest struct {
	sdk.CertificateRequest
}

func NewRequestRequest() *RequestRequest {
	return &RequestRequest{}
}
