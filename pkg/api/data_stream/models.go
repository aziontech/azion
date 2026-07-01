package datastream

import (
	"encoding/json"
	"time"

	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/data-stream-api"
)

type CreateRequest struct {
	sdk.DataStreamRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchedDataStreamRequest
}

func NewUpdateRequest() *UpdateRequest {
	return &UpdateRequest{}
}

// Response, ListResponse and DataStream are decoded by the CLI instead of relying
// on the generated SDK types: the SDK fails to unmarshal the polymorphic
// inputs/transform/outputs returned by the Data Stream API (oneOf discriminator
// mismatch), which would otherwise leave the whole payload empty. The polymorphic
// fields are kept as raw JSON since the CLI only displays the scalar attributes.
type Response struct {
	State string     `json:"state"`
	Data  DataStream `json:"data"`
}

type ListResponse struct {
	Count   int64        `json:"count"`
	Results []DataStream `json:"results"`
}

type DataStream struct {
	Id             int64           `json:"id"`
	Name           string          `json:"name"`
	Active         bool            `json:"active"`
	LastEditor     string          `json:"last_editor"`
	ProductVersion string          `json:"product_version"`
	Created        time.Time       `json:"created"`
	LastModified   time.Time       `json:"last_modified"`
	Inputs         json.RawMessage `json:"inputs,omitempty"`
	Transform      json.RawMessage `json:"transform,omitempty"`
	Outputs        json.RawMessage `json:"outputs,omitempty"`
}

func (d DataStream) GetId() int64               { return d.Id }
func (d DataStream) GetName() string            { return d.Name }
func (d DataStream) GetActive() bool            { return d.Active }
func (d DataStream) GetLastEditor() string      { return d.LastEditor }
func (d DataStream) GetLastModified() time.Time { return d.LastModified }
func (d DataStream) GetProductVersion() string  { return d.ProductVersion }
