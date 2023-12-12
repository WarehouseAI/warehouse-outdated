package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type AuthScheme string

const (
	Bearer AuthScheme = "Bearer"
	ApiKey AuthScheme = "ApiKey"
	Basic  AuthScheme = "Basic"
)

// TODO: Подумать над синхронизацией с сервисом пользователей.
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
