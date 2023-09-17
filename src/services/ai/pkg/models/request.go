package models

import dbm "warehouse/src/internal/db/models"

type (
	CreateAIRequest struct {
		Name       string         `json:"name"`
		AuthScheme dbm.AuthScheme `json:"authScheme"`
	}
)
