package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type AuthScheme string
type RequestType string
type PayloadType string

const (
	Bearer AuthScheme = "Bearer"
	ApiKey AuthScheme = "ApiKey"
	Basic  AuthScheme = "Basic"
)

const (
	Post    RequestType = "POST"
	Get     RequestType = "GET"
	Put     RequestType = "PUT"
	Update  RequestType = "UPDATE"
	Delete  RequestType = "DELETE"
	Patch   RequestType = "PATCH"
	Head    RequestType = "HEAD"
	Options RequestType = "OPTIONS"
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
		ID          uuid.UUID              `json:"id" gorm:"type:uuid;primarykey"`
		AI          uuid.UUID              `json:"ai" gorm:"type:uuid"`
		Name        string                 `json:"name" gorm:"type:string"`
		Payload     map[string]interface{} `json:"payload" gorm:"type:json;not null"`
		PayloadType PayloadType            `json:"payload_type" gorm:"type:PayloadType;not null"`
		RequestType RequestType            `json:"request_type" gorm:"type:RequestType;not null"`
		URL         string                 `json:"url" gorm:"type:string;unique;not null"`
		CreatedAt   time.Time              `json:"created_at" gorm:"type:time"`
		UpdateAt    time.Time              `json:"updated_at" gorm:"type:time"`
	}
)
