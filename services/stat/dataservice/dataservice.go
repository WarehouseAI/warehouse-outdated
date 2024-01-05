package dataservice

import (
	e "warehouseai/stat/errors"

	"github.com/gofrs/uuid"
)

type UsersStatInterface interface {
	GetNumOfUsers() (uint, *e.DBError)
	GetNumOfDevelopers() (uint, *e.DBError)
}

type AiStatInterface interface {
	GetNumOfAiUses(id uuid.UUID) (uint, *e.DBError)
}
