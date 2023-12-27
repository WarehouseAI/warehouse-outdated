package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
)

type AuthScheme string
type RequestScheme string
type IOType string
type PayloadType string

const (
	Bearer AuthScheme = "Bearer"
	ApiKey AuthScheme = "ApiKey"
	Basic  AuthScheme = "Basic"
)
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

type AI struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`
	Owner         uuid.UUID  `json:"owner" gorm:"type:uuid;not null"`
	Commands      []Command  `json:"commands" gorm:"foreignKey:AIID"`
	Description   string     `json:"description" gorm:"type:string;not null"`
	BackgroundUrl string     `json:"background_url" gorm:"type:string;not null"`
	Name          string     `json:"name" gorm:"type:string;unique;not null"`
	ApiKey        string     `json:"-" gorm:"type:string;not null"`
	AuthScheme    AuthScheme `json:"-" gorm:"type:AuthScheme;not null"`
	Used          int        `json:"used" gorm:"type:int;default:0"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:time"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:time"`
}

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
