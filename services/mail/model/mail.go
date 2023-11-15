package model

type EmailReceivedEvent struct {
	Data Email `json:"data"`
}

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
