package main

import (
	"errors"
	"gopkg.in/mailgun/mailgun-go.v1"
	"log"
	"os"
)

type Message struct {
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	ReplyTo  string                 `json:"reply_to"`
	Subject  string                 `json:"subject"`
	Text     string                 `json:"text"`
	Html     string                 `json:"html"`
	Track    bool                   `json:"track"`
	Template string                 `json:template`
	Params   map[string]interface{} `json:params`
}

type MailService struct {
	Service         mailgun.Mailgun
	TestMode        bool
	Messages        chan *Message
	templateManager *TemplateManager
}

func NewMailService() *MailService {
	domain := os.Getenv("DOMAIN")
	apiKey := os.Getenv("API_KEY")
	publicApiKey := os.Getenv("PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	testMode := os.Getenv("TEST_MODE")
	service := &MailService{Service: mg, TestMode: "" != testMode, templateManager: NewTemplateManager()}
	service.run()
	return service
}

func (ms *MailService) run() {
	ms.Messages = make(chan *Message, 50)
	go func() {
		for {
			msg := <-ms.Messages
			if "" != msg.Template {
				ms.SendTemplatedEmail(msg)
				continue
			}
			ms.SendEmail(msg)
		}
	}()
}

func (ms *MailService) QueueEmail(message *Message) {
	ms.Messages <- message
}

func (ms *MailService) SendEmail(msg *Message) error {
	if msg.To == "" {
		log.Printf("Error: To address is required for sending emails\n")
		return errors.New("missing param: to address not set")
	}
	message := ms.Service.NewMessage(
		msg.From,
		msg.Subject,
		msg.Text,
		msg.To)
	if "" != msg.Html {
		message.SetHtml(msg.Html)
	}
	if "" != msg.ReplyTo {
		message.SetReplyTo(msg.ReplyTo)
	}
	message.SetTracking(msg.Track)
	return ms.send(message)
}

func (ms *MailService) SendTemplatedEmail(msg *Message) error {
	if msg.To == "" {
		log.Printf("Error: To address is required for sending emails\n")
		return errors.New("missing param: to address not set")
	}
	err := ms.assembleTemplate(msg)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return err
	}

	if ms.TestMode {
		log.Printf("----------\nSubject: %s\nText: %s\nHtml: %s\n----------\n", msg.Subject, msg.Text, msg.Html)
	}
	message := ms.Service.NewMessage(
		msg.From,
		msg.Subject,
		msg.Text,
		msg.To)

	if "" != msg.Html {
		message.SetHtml(msg.Html)
	}
	if "" != msg.ReplyTo {
		message.SetReplyTo(msg.ReplyTo)
	}
	message.SetTracking(msg.Track)

	return ms.send(message)
}

func (ms *MailService) assembleTemplate(msg *Message) error {
	emailTpl, err := ms.templateManager.ImportTemplate(msg.Template)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return err
	}

	if emailTpl.From != "" {
		msg.From = emailTpl.From
	}

	if emailTpl.Subject != "" {
		msg.Subject = emailTpl.Subject
	}

	text, err := AssembleTemplate(emailTpl.Text, msg)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return err
	}
	msg.Text = text

	if nil != emailTpl.Html {
		html, err := AssembleTemplate(emailTpl.Html, msg)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return err
		}
		msg.Html = html
	}

	return nil
}

func (ms *MailService) send(message *mailgun.Message) error {
	if ms.TestMode {
		message.EnableTestMode()
	}

	resp, id, err := ms.Service.Send(message)

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return err
	}
	log.Printf("ID: %s Resp: %s\n", id, resp)
	return err
}
