package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type (
	User struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
		Username  string    `json:"name" gorm:"type:string;unique"`
		Picture   string    `json:"picture" gorm:"type:string"`
		Password  string    `json:"-" gorm:"type:string;not null"`
		Email     string    `json:"email" gorm:"type:string;not null;unique"`
		ViaGoogle bool      `json:"via_google" gorm:"default:false;not null"`
		OwnedAi   []AI      `json:"owned_ai" gorm:"foreignKey:Owner"`
		CreatedAt time.Time `json:"created_at" gorm:"type:time"`
		UpdateAt  time.Time `json:"updated_at" gorm:"type:time"`
	}
)
