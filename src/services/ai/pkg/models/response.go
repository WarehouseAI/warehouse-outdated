package models

import dbm "warehouse/src/internal/db/models"

type (
	CreateAIResponse struct {
		Name       string         `json:"name"`
		AuthScheme dbm.AuthScheme `json:"authScheme"`
		ApiKey     string         `json:"apiKey"`
	}
)
