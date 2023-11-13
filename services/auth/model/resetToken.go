package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type ResetToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
	UserId    uuid.UUID `json:"-" gorm:"type:uuid;not null;unique"`
	Token     string    `json:"-" gorm:"type:string;not null"`
	ExpiresAt time.Time `json:"-" gorm:"type:time"`
	CreatedAt time.Time `json:"-" gorm:"type:time"`
}
