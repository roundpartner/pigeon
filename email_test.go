package main

import (
	"os"
	"testing"
)

var FromEmail string
var ToEmail string

func TestInit(t *testing.T) {
	FromEmail = os.Getenv("FROM_EMAIL")
	ToEmail = os.Getenv("TO_EMAIL")
	if ToEmail == "" {
		t.Fail()
	}
	if FromEmail == "" {
		t.Fail()
	}
}

func TestMessageDefaults(t *testing.T) {
	msg := &Message{}
	if msg.Html != "" {
		t.Errorf("Html was not false: %s", msg.Html)
	}
	if msg.Track != false {
		t.Errorf("Track was not false: %t", msg.Track)
	}
}

func TestSendsEmail(t *testing.T) {
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, Subject: "Queued Message", Text: "This tests that messages can be queued"}
	err := service.SendEmail(&message)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestSendEmailWithReport(t *testing.T) {
	service := NewMailService()
	message := Message{
		From:    FromEmail,
		To:      ToEmail,
		Subject: "Queued Message",
		Text:    "This tests that messages can be queued",
		Report:  true,
	}
	err := service.SendEmail(&message)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
	if message.Subject != "Queued Message [Spam: false Score: 1.000000]" {
		t.Errorf("Subject %s did not match", message.Subject)
		t.FailNow()
	}
}

func TestFromEmailBlocked(t *testing.T) {
	os.Setenv("BLACK_LISTED_ADDRESSES", `\.co\.uk$`)
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, Subject: "Blocked Message", Text: "This tests that messages can be blocked"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLACK_LISTED_ADDRESSES")
	if err == nil {
		t.FailNow()
	}
	if "black listed email" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestReplyToEmailBlocked(t *testing.T) {
	os.Setenv("BLACK_LISTED_ADDRESSES", `mailinator\.com$`)
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, ReplyTo: "test@mailinator.com", Subject: "Blocked Message", Text: "This tests that messages can be blocked"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLACK_LISTED_ADDRESSES")
	if err == nil {
		t.FailNow()
	}
	if "black listed email" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestContentEmailBlocked(t *testing.T) {
	os.Setenv("BLACK_LISTED_CONTENT", `blocked|test`)
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, ReplyTo: "test@mailinator.com", Subject: "Blocked Message", Text: "This tests that messages can be blocked when keywords are being filtered"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLACK_LISTED_CONTENT")
	if err == nil {
		t.FailNow()
	}
	if "black listed phrase" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestQueuesEmail(t *testing.T) {
	service := NewMailService()

	if nil == service.Messages {
		t.Error("Message queue has not been set")
		t.FailNow()
	}

	message := Message{From: FromEmail, To: ToEmail, Subject: "Queued Message", Text: "This tests that messages can be queued"}
	service.QueueEmail(&message)
}

func TestSendTemplatedEmail(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, Subject: "Queued Message", Template: "test"}
	err := service.SendTemplatedEmail(&message)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
		t.FailNow()
	}
}

func TestSendTemplatedEmailWithReport(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, Subject: "Queued Message", Template: "test"}
	message.Report = true
	err := service.SendTemplatedEmail(&message)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
		t.FailNow()
	}
	if message.Subject != "Queued Message [Spam: false Score: 1.000000]" {
		t.Errorf("Subject %s did not match", message.Subject)
		t.FailNow()
	}
}

func TestAssembleTemplate(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, Template: "test"}
	service.AssembleTemplate(&message)
	if message.From == "" {
		t.Errorf("Error: from for email was not assembled\n")
		t.FailNow()
	}
	if message.Subject == "" {
		t.Errorf("Error: subject for email was not assembled\n")
		t.FailNow()
	}
	if message.Text == "" {
		t.Errorf("Error: text for email was not assembled\n")
		t.FailNow()
	}
	if message.Html == "" {
		t.Errorf("Error: html for email was not assembled\n")
		t.FailNow()
	}
}

func TestAssembleTemplateDoesNotChangeSubject(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, Subject: "Queued Message", Template: "test"}
	service.AssembleTemplate(&message)
	if message.Subject != "Queued Message" {
		t.Errorf("Error: subject for email was not assembled\n")
		t.FailNow()
	}
}
