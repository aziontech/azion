package crl

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.CertificateRevocationList
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{CertificateRevocationList: *sdk.NewCertificateRevocationListWithDefaults()}
}

type UpdateRequest struct {
	sdk.PatchedCertificateRevocationList
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
