package main

import (
	"gopkg.in/mailgun/mailgun-go.v1"
	"log"
	"os"
)

type MailService struct {
	Service  mailgun.Mailgun
	TestMode bool
}

func NewMailService() *MailService {
	domain := os.Getenv("DOMAIN")
	apiKey := os.Getenv("API_KEY")
	publicApiKey := os.Getenv("PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	testMode := os.Getenv("TEST_MODE")
	service := &MailService{Service: mg, TestMode: "" != testMode}
	return service
}

func (ms MailService) SendEmail(from string, to string, subject string, text string) error {
	message := ms.Service.NewMessage(
		from,
		subject,
		text,
		to)
	if ms.TestMode {
		message.EnableTestMode()
	}
	resp, id, err := ms.Service.Send(message)

	if err != nil {
		return err
	}
	log.Printf("ID: %s Resp: %s\n", id, resp)
	return err
}
