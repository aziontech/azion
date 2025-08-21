package edgefunction

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
)

type CreateRequest struct {
	sdk.EdgeFunctionsRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedEdgeFunctionsRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
