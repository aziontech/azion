package edgefunction

import (
	sdk "github.com/aziontech/azionapi-go-sdk/edgefunctions"
)

type CreateRequest struct {
	sdk.CreateEdgeFunctionRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchEdgeFunctionRequest
	Id int64
}

func NewUpdateRequest(id int64) *UpdateRequest {
	return &UpdateRequest{Id: id}
}

type EdgeFunctionResponse interface {
	GetId() int64
	GetName() string
	GetActive() bool
	GetLanguage() string
	GetReferenceCount() int64
	GetModified() string
	GetInitiatorType() string
	GetLastEditor() string
	GetFunctionToRun() string
	GetJsonArgs() interface{}
	GetCode() string
}
