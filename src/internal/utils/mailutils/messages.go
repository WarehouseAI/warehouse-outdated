package mailutils

import (
	"fmt"
	"os"
)

type EmailType string

type Email struct {
	From    string
	To      string
	Subject string
	Message string
}

const (
	EmailVerify EmailType = "email_verify"
)

func NewMessage(emailType EmailType, to string, link string) Email {
	from := fmt.Sprintf("%s@%s", os.Getenv("MAIL_USER"), os.Getenv("MAIL_DOMAIN"))

	messages := map[EmailType]Email{
		EmailVerify: {
			From:    from,
			To:      to,
			Subject: "Подтверждение электронной почты",
			Message: fmt.Sprintf("Для подтверждения электронной почты пройдите по ссылке:\n%s", link),
		},
	}

	return messages[emailType]
}
