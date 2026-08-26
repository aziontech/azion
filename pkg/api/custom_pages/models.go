package custompages

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.CustomPageRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedCustomPageRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
