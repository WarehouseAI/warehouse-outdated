package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type UserRole string

const (
	Developer UserRole = "DEVELOPER"
	Base      UserRole = "BASE"
)

type User struct {
	ID         uuid.UUID   `json:"id" gorm:"type:uuid;primarykey"`
	Firstname  string      `json:"firstname" gorm:"type:string;not null"`
	Lastname   string      `json:"lastname" gorm:"type:string;not null"`
	Username   string      `json:"username" gorm:"type:string;unique"`
	Picture    string      `json:"picture" gorm:"type:string"`
	Password   string      `json:"-" gorm:"type:string;not null"`
	Email      string      `json:"email" gorm:"type:string;not null;unique"`
	ViaGoogle  bool        `json:"via_google;omitempty" gorm:"default:false;not null"`
	Verified   bool        `json:"-" gorm:"default:false;not null"`
	Role       UserRole    `json:"role" gorm:"type:UserRole;default:Base;not null"`
	FavoriteAi []uuid.UUID `json:"favorite_ai" gorm:"type:uuid"`
	OwnedAi    []uuid.UUID `json:"owned_ai" gorm:"type:uuid"`
	CreatedAt  time.Time   `json:"created_at" gorm:"type:time"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"type:time"`
}

type UserFavorites struct {
	ID     uint      `json:"-" gorm:"primarykey"`
	AiId   uuid.UUID `json:"ai_id" gorm:"type:uuid;not null"`
	UserId uuid.UUID `json:"user_id" gorm:"uuid"`
}
