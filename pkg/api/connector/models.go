package connector

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.ConnectorRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedConnectorRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
