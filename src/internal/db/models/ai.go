package models

import (
	"time"

	"github.com/gofrs/uuid"
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

const (
	FormData PayloadType = "FormData"
	Json     PayloadType = "JSON"
)

type (
	All interface {
		AI | Command
	}

	AI struct {
		ID         uuid.UUID  `json:"id" gorm:"type:uuid;primarykey"`
		Owner      uuid.UUID  `json:"owner" gorm:"type:uuid"`
		Commands   []Command  `json:"commands" gorm:"foreignKey:AI"`
		Name       string     `json:"name" gorm:"type:string;unique;not null"`
		ApiKey     string     `json:"-" gorm:"type:string;not null"`
		AuthScheme AuthScheme `json:"auth_scheme" gorm:"type:AuthScheme;not null"`
		CreatedAt  time.Time  `json:"created_at" gorm:"type:time"`
		UpdateAt   time.Time  `json:"updated_at" gorm:"type:time"`
	}

	Command struct {
		ID            uuid.UUID              `json:"id" gorm:"type:uuid;primarykey"`
		AI            uuid.UUID              `json:"ai" gorm:"type:uuid"`
		Name          string                 `json:"name" gorm:"type:string"`
		Payload       map[string]interface{} `json:"payload" gorm:"type:json;not null"`
		PayloadType   PayloadType            `json:"payload_type" gorm:"type:PayloadType;not null"`
		RequestScheme RequestScheme          `json:"request_type" gorm:"type:RequestType;not null"`
		InputType     IOType                 `json:"input_type" gorm:"type:IOType;not null"`
		OutputType    IOType                 `json:"output_type" gorm:"type:IOType;not null"`
		URL           string                 `json:"url" gorm:"type:string;unique;not null"`
		CreatedAt     time.Time              `json:"created_at" gorm:"type:time"`
		UpdateAt      time.Time              `json:"updated_at" gorm:"type:time"`
	}
)
