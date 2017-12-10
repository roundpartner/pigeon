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
	err := service.SendEmail(FromEmail, ToEmail, "Test Subject", "This is a really cool message")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}
