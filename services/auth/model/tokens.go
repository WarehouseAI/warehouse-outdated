package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Tokens interface {
	ResetToken | VerificationToken
}

type ResetToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`
	UserId    uuid.UUID `json:"-" gorm:"type:uuid;not null;unique"`
	Token     string    `json:"-" gorm:"type:string;not null"`
	ExpiresAt time.Time `json:"-" gorm:"type:time"`
	CreatedAt time.Time `json:"-" gorm:"type:time"`
}

type VerificationToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`
	UserId    string    `json:"-" gorm:"type:uuid;not null;unique"`
	Token     string    `json:"-" gorm:"type:string;not null"`
	ExpiresAt time.Time `json:"-" gorm:"type:time"`
	CreatedAt time.Time `json:"-" gorm:"type:time"`
}
