package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type (
	User struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid"`
		Username  string    `json:"name" gorm:"type:string;unique"`
		Picture   string    `json:"picture" gorm:"type:string"`
		Password  string    `json:"-" gorm:"type:string"`
		Email     string    `json:"email" gorm:"type:string;not null;unique"`
		ViaGoogle bool      `json:"viaGoogle" gorm:"default:false;not null"`
		CreatedAt time.Time `json:"createdAt" gorm:"type:time"`
		UpdateAt  time.Time `json:"updatedAt" gorm:"type:time"`
	}
)
