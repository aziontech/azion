package workloads

import (
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
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
	GetDomains() []string
	GetLastEditor() string
	GetActive() bool
}

type DeploymentResponse interface {
	GetId() int64
	GetCurrent() bool
	GetActive() bool
	GetName() string
}
