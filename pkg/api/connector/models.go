package connector

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.ConnectorRequest2
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedConnectorRequest2
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
