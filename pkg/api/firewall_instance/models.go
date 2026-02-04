package firewallinstance

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

type CreateRequest struct {
	sdk.FirewallFunctionInstanceRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedFirewallFunctionInstanceRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
