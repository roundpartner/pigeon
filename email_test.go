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
}

func TestMessageDefaults(t *testing.T) {
	msg := &Message{}
	if msg.Html != "" {
		t.Errorf("Html was not false: %s", msg.Html)
	}
	if msg.Track != false {
		t.Errorf("Track was not false: %s", msg.Track)
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
	message := Message{From: FromEmail, To: ToEmail, Subject: "Queued Message", Template: "This is a test template"}
	err := service.SendTemplatedEmail(&message)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}
