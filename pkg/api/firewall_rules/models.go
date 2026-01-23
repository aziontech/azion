package firewallrules

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
)

type CreateRequest struct {
	sdk.FirewallRuleRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedFirewallRuleRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}
