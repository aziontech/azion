package firewall

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
)

type CreateRequest struct {
	sdk.FirewallRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedFirewallRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
