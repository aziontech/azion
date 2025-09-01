package edgeconnector

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
)

type CreateRequest struct {
	sdk.ConnectorPolymorphicRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedConnectorPolymorphicRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
