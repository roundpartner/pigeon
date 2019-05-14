package main

import (
	"errors"
	"gopkg.in/mailgun/mailgun-go.v1"
	"log"
	"os"
	"regexp"
)

type Message struct {
	From      string                 `json:"from"`
	FromName  string                 `json:"from_name,omitempty"`
	To        string                 `json:"to"`
	ReplyTo   string                 `json:"reply_to"`
	Subject   string                 `json:"subject"`
	Text      string                 `json:"text"`
	Html      string                 `json:"html"`
	Track     bool                   `json:"track"`
	Template  string                 `json:"template"`
	SenderIp  string                 `json:"sender_ip,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Website   string                 `json:"website,omitempty"`
	Params    map[string]interface{} `json:"params"`
	Report    bool                   `json:"report,omitempty"`
	IsSpam    bool
	SpamScore float64
}

type MailService struct {
	Service            mailgun.Mailgun
	TestMode           bool
	Messages           chan *Message
	templateManager    *TemplateManager
	BlackListedAddress *regexp.Regexp
	BlackListedContent *regexp.Regexp
}

func NewMailService() *MailService {
	domain := os.Getenv("DOMAIN")
	apiKey := os.Getenv("API_KEY")
	publicApiKey := os.Getenv("PUBLIC_API_KEY")
	os.Setenv("MG_DOMAIN", domain)
	os.Setenv("MG_API_KEY", apiKey)
	os.Setenv("MG_PUBLIC_API_KEY", publicApiKey)

	mg, err := mailgun.NewMailgunFromEnv()

	url := os.Getenv("MG_URL")
	if url == "" {
		mg.SetAPIBase(mailgun.ApiBase)
	}
	if err != nil {
		log.Printf("[INFO] [%s] %s", ServiceName, err.Error())
		os.Exit(1)
	}
	testMode := os.Getenv("TEST_MODE")
	service := &MailService{Service: mg, TestMode: "" != testMode, templateManager: NewTemplateManager()}
	blackListedAddress, isSet := os.LookupEnv("BLACK_LISTED_ADDRESSES")
	if isSet && "" != blackListedAddress {
		service.BlackListedAddress = regexp.MustCompile(blackListedAddress)
	}
	blackListedContent, isSet := os.LookupEnv("BLACK_LISTED_CONTENT")
	if isSet && "" != blackListedContent {
		service.BlackListedContent = regexp.MustCompile("(?i)" + blackListedContent)
	}
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
		log.Printf("[ERROR] [%s] To address is required for sending emails", ServiceName)
		return errors.New("missing param: to address not set")
	}
	if nil != ms.BlackListedAddress && ms.BlackListedAddress.MatchString(msg.From) {
		log.Printf("[INFO] [%s] From address has been blacklisted\n", ServiceName)
		return errors.New("black listed email")
	}
	if nil != ms.BlackListedAddress && ms.BlackListedAddress.MatchString(msg.ReplyTo) {
		log.Printf("[INFO] [%s] ReplyTo address has been blacklisted", ServiceName)
		return errors.New("black listed email")
	}
	if nil != ms.BlackListedAddress && ms.BlackListedAddress.MatchString(msg.To) {
		log.Printf("[INFO] [%s] From address has been blacklisted", ServiceName)
		return errors.New("black listed sender email")
	}
	if nil != ms.BlackListedContent {
		if ms.BlackListedContent.MatchString(msg.Text) {
			log.Printf("[INFO] [%s] Text has been blacklisted", ServiceName)
			return errors.New("black listed phrase")
		}
	}
	if msg.SenderIp != "" && CheckBlackList(msg.SenderIp) {
		log.Printf("[INFO] [%s] sender ip has been blacklisted", ServiceName)
		return errors.New("black listed ip")
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
	log.Printf("[INFO] [%s] Sending email", ServiceName)
	return ms.send(message)
}

func (ms *MailService) SendTemplatedEmail(msg *Message) error {
	if msg.To == "" {
		log.Printf("[ERROR] [%s] To address is required for sending emails", ServiceName)
		return errors.New("missing param: to address not set")
	}
	err := ms.AssembleTemplate(msg)
	if err != nil {
		log.Printf("[ERROR] [%s] %s\n", ServiceName, err.Error())
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

func (ms *MailService) AssembleTemplate(msg *Message) error {
	emailTpl, err := ms.templateManager.ImportTemplate(msg.Template)
	if err != nil {
		log.Printf("[ERROR] [%s] %s\n", ServiceName, err.Error())
		return err
	}

	if emailTpl.From != "" {
		msg.From = emailTpl.From
	}

	if msg.Subject == "" && emailTpl.Subject != "" {
		msg.Subject = emailTpl.Subject
	}

	text, err := AssembleTemplate(emailTpl.Text, msg)
	if err != nil {
		log.Printf("[ERROR] [%s] %s\n", ServiceName, err.Error())
		return err
	}
	msg.Text = text

	if nil != emailTpl.Html {
		html, err := AssembleTemplate(emailTpl.Html, msg)
		if err != nil {
			log.Printf("[ERROR] [%s] %s\n", ServiceName, err.Error())
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
		log.Printf("[ERROR] [%s] %s\n", ServiceName, err.Error())
		return err
	}
	log.Printf("[ERROR] [%s] ID: %s Resp: %s\n", ServiceName, id, resp)
	return err
}
