package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
)

type PayloadType string

const (
	FormData PayloadType = "FormData"
	Json     PayloadType = "JSON"
)

type RequestScheme string

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

type IOType string

const (
	Image IOType = "Image"
	Audio IOType = "Audio"
	Text  IOType = "Text"
)

type FieldClass string

const (
	Permanent FieldClass = "permanent"
	Optional  FieldClass = "optional"
	Free      FieldClass = "free"
)

type DataType string

const (
	String DataType = "string"
	Number DataType = "number"
	File   DataType = "file"
	Bool   DataType = "bool"
	Object DataType = "object"
)

type AiCommand struct {
	ID          uuid.UUID         `json:"id" gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`
	AIID        uuid.UUID         `json:"ai_id" gorm:"type:uuid"`
	Name        string            `json:"name" gorm:"type:string"`
	Payload     datatypes.JSONMap `json:"payload" gorm:"type:json;not null"`
	PayloadType string            `json:"payload_type" gorm:"type:string;not null"`
	RequestType string            `json:"request_type" gorm:"type:string;not null"`
	InputType   string            `json:"input_type" gorm:"type:string;not null"`
	OutputType  string            `json:"output_type" gorm:"type:string;not null"`
	URL         string            `json:"url" gorm:"type:string;unique;not null"`
	CreatedAt   time.Time         `json:"created_at" gorm:"type:time"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"type:time"`
}

type CommandFieldParams struct {
	Class    FieldClass    `json:"class"`
	Values   []interface{} `json:"values"`
	DataType DataType      `json:"data_type"`
}
