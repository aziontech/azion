package waf

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.WAFRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedWAFRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
