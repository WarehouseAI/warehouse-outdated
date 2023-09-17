package models

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

type (
	AI struct {
		ID         uuid.UUID  `json:"id" gorm:"type:uuid;primarykey"`
		Owner      uuid.UUID  `json:"owner" gorm:"type:uuid"`
		Name       string     `json:"name" gorm:"type:string;unique;not null"`
		ApiKey     string     `json:"-" gorm:"type:string;not null"`
		AuthScheme AuthScheme `json:"authScheme" gorm:"type:AuthScheme;not null"`
		CreatedAt  time.Time  `json:"createdAt" gorm:"type:time"`
		UpdateAt   time.Time  `json:"updatedAt" gorm:"type:time"`
	}
)
