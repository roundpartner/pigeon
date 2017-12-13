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
