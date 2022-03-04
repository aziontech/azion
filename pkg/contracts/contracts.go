package contracts

import (
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
)

type ListOptions struct {
	Details  bool
	OrderBy  string
	Sort     string
	Page     int64
	PageSize int64
	Filter   string
}

type DescribeOptions struct {
	OutPath string
	Format  string
}

type UpdateRequestResource struct {
	sdk.UpdateResourceRequest
	Id int64
}

type UpdateRequestService struct {
	sdk.UpdateServiceRequest
	Id int64
}
