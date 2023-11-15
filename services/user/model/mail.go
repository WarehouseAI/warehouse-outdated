package model

type EmailRequest struct {
	Data Email `json:"data"`
}

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
