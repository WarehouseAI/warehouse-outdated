package models

import dbm "warehouse/src/internal/db/models"

type (
	CreateAIResponse struct {
		Name       string         `json:"name"`
		AuthScheme dbm.AuthScheme `json:"auth_scheme"`
		ApiKey     string         `json:"api_key"`
	}
)
