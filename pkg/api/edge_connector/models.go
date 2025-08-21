package edgeconnector

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
)

type CreateRequest struct {
	sdk.EdgeConnectorPolymorphicRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedEdgeConnectorPolymorphicRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
