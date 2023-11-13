package adapter

import (
	"warehouseai/auth/model"
)

type MailProducerInterface interface {
	SendEmail(email model.Email) error
}
