package workloads

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
)

type Client struct {
	apiClient *sdk.APIClient
}

type CreateRequest struct {
	sdk.WorkloadRequest
}

type UpdateRequest struct {
	sdk.PatchedWorkloadRequest
	Id int64
}

type WorkloadResponse interface {
	GetId() int64
	GetName() string
	GetDomains() []sdk.DomainInfo
	GetLastEditor() string
	GetActive() bool
	GetAlternateDomains() []string
}

type DeploymentResponse interface {
	GetId() int64
	GetTag() string
	GetCurrent() bool
	GetBinds() sdk.WorkloadDeploymentBinds
}
