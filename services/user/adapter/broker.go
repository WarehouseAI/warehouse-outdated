package adapter

import "warehouseai/user/model"

type MailProducerInterface interface {
	SendEmail(email model.Email) error
}
