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

type Command struct {
	ID            uuid.UUID         `json:"id" gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`
	AIID          uuid.UUID         `json:"ai_id" gorm:"type:uuid"`
	Name          string            `json:"name" gorm:"type:string"`
	Payload       datatypes.JSONMap `json:"payload" gorm:"type:json;not null"`
	PayloadType   PayloadType       `json:"payload_type" gorm:"type:PayloadType;not null"`
	RequestScheme RequestScheme     `json:"request_type" gorm:"type:RequestScheme;not null"`
	InputType     IOType            `json:"input_type" gorm:"type:IOType;not null"`
	OutputType    IOType            `json:"output_type" gorm:"type:IOType;not null"`
	URL           string            `json:"url" gorm:"type:string;unique;not null"`
	CreatedAt     time.Time         `json:"created_at" gorm:"type:time"`
	UpdatedAt     time.Time         `json:"updated_at" gorm:"type:time"`
}

type CommandFieldParams struct {
	Class    FieldClass    `json:"class"`
	Values   []interface{} `json:"values"`
	DataType DataType      `json:"data_type"`
}
