package variables

import (
	"time"

	sdk "github.com/aziontech/azionapi-go-sdk/variables"
)

type Request struct {
	sdk.VariableCreate
	Uuid string
}

type Response interface {
	GetUuid() string
	GetKey() string
	GetValue() string
	GetSecret() bool
	GetLastEditor() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}
