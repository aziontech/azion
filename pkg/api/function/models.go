package function

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.FunctionsRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedFunctionsRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
