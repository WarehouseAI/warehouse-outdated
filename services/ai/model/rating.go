package model

import "github.com/gofrs/uuid"

type RatingPerUser struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`
	UserId uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	AiId   uuid.UUID `json:"ai_id" gorm:"type:uuid;not null"`
	Rate   int16     `json:"rate" gorm:"check:rate > 0;check:rate <= 5;not null"`
}
