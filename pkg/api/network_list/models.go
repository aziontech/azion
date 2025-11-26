package networklist

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
)

type CreateRequest struct {
	sdk.NetworkListDetailRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedNetworkListDetailRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
