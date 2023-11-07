package postgresdb

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
)

type All interface {
	AI | Command | User | ResetToken
}

type AuthScheme string
type RequestScheme string
type IOType string
type PayloadType string
type UserRole string

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

const (
	Developer UserRole = "DEVELOPER"
	Base      UserRole = "BASE"
)

type (
	AI struct {
		ID          uuid.UUID  `json:"id" gorm:"type:uuid;primarykey"`
		Owner       uuid.UUID  `json:"owner" gorm:"foreignKey:ID"`
		FavoriteFor []*User    `json:"-" gorm:"many2many:user_favorites"`
		Commands    []Command  `json:"commands" gorm:"foreignKey:AIID"`
		Description string     `json:"description" gorm:"type:string;not null"`
		Name        string     `json:"name" gorm:"type:string;unique;not null"`
		ApiKey      string     `json:"-" gorm:"type:string;not null"`
		AuthScheme  AuthScheme `json:"-" gorm:"type:AuthScheme;not null"`
		Used        int        `json:"used" gorm:"type:int;default:0"`
		CreatedAt   time.Time  `json:"created_at" gorm:"type:time"`
		UpdatedAt   time.Time  `json:"updated_at" gorm:"type:time"`
	}

	Command struct {
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

	User struct {
		ID               uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
		Firstname        string    `json:"firstname" gorm:"type:string;not null"`
		Lastname         string    `json:"lastname" gorm:"type:string;not null"`
		Username         string    `json:"username" gorm:"type:string;unique"`
		Picture          string    `json:"picture" gorm:"type:string"`
		Password         string    `json:"-" gorm:"type:string;not null"`
		Email            string    `json:"email" gorm:"type:string;not null;unique"`
		ViaGoogle        bool      `json:"via_google;omitempty" gorm:"default:false;not null"`
		Verified         bool      `json:"-" gorm:"default:false;not null"`
		VerificationCode *string   `json:"-" gorm:"type:string"`
		Role             UserRole  `json:"role" gorm:"type:UserRole;default:Base;not null"`
		FavoriteAi       []*AI     `json:"favorite_ai" gorm:"many2many:user_favorites"`
		OwnedAi          []AI      `json:"owned_ai" gorm:"foreignKey:Owner"`
		CreatedAt        time.Time `json:"created_at" gorm:"type:time"`
		UpdatedAt        time.Time `json:"updated_at" gorm:"type:time"`
	}

	ResetToken struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
		UserId    uuid.UUID `json:"-" gorm:"type:uuid;not null;unique"`
		Token     string    `json:"-" gorm:"type:string;not null"`
		ExpiresAt time.Time `json:"-" gorm:"type:time"`
		CreatedAt time.Time `json:"-" gorm:"type:time"`
	}
)
