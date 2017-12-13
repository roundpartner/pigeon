package main

import (
	"gopkg.in/mailgun/mailgun-go.v1"
	"log"
	"os"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type MailService struct {
	Service  mailgun.Mailgun
	TestMode bool
	Messages chan *Message
}

func NewMailService() *MailService {
	domain := os.Getenv("DOMAIN")
	apiKey := os.Getenv("API_KEY")
	publicApiKey := os.Getenv("PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	testMode := os.Getenv("TEST_MODE")
	service := &MailService{Service: mg, TestMode: "" != testMode}
	service.run()
	return service
}

func (ms *MailService) run() {
	ms.Messages = make(chan *Message, 50)
	go func() {
		for {
			msg := <-ms.Messages
			ms.SendEmail(msg)
		}
	}()
}

func (ms *MailService) QueueEmail(message *Message) {
	ms.Messages <- message
}

func (ms *MailService) SendEmail(msg *Message) error {
	message := ms.Service.NewMessage(
		msg.From,
		msg.Subject,
		msg.Text,
		msg.To)
	if ms.TestMode {
		message.EnableTestMode()
	}
	message.SetTracking(false)
	resp, id, err := ms.Service.Send(message)

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return err
	}
	log.Printf("ID: %s Resp: %s\n", id, resp)
	return err
}
