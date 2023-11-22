package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
)

type RequestScheme string
type IOType string
type PayloadType string

const (
	FormData PayloadType = "FormData"
	Json     PayloadType = "JSON"
)

const (
	Post    RequestScheme = "POST"
	Get     RequestScheme = "GET"
	Put     RequestScheme = "PUT"
	Update  RequestScheme = "UPDATE"
	Delete  RequestScheme = "DELETE"
	Patch   RequestScheme = "PATCH"
	Head    RequestScheme = "HEAD"
	Options RequestScheme = "OPTIONS"
)

const (
	Image IOType = "Image"
	Text  IOType = "Text"
)

type Command struct {
	ID            uuid.UUID         `json:"id" gorm:"type:uuid;primarykey"`
	AIID          uuid.UUID         `json:"ai_id" gorm:"type:uuid"`
	Name          string            `json:"name" gorm:"type:string"`
	Payload       datatypes.JSONMap `json:"payload" gorm:"type:json;not null"`
	PayloadType   PayloadType       `json:"payload_type" gorm:"type:PayloadType;not null"`
	RequestScheme RequestScheme     `json:"request_type" gorm:"type:RequestType;not null"`
	InputType     IOType            `json:"input_type" gorm:"type:IOType;not null"`
	OutputType    IOType            `json:"output_type" gorm:"type:IOType;not null"`
	URL           string            `json:"url" gorm:"type:string;unique;not null"`
	CreatedAt     time.Time         `json:"created_at" gorm:"type:time"`
	UpdatedAt     time.Time         `json:"updated_at" gorm:"type:time"`
}
