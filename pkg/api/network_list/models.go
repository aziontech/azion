package networklist

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.NetworkListRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedNetworkListRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
