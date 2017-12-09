package main

import (
	"testing"
	"os"
)

var FromEmail string
var ToEmail string

func TestInit(t *testing.T) {
	FromEmail = os.Getenv("FROM_EMAIL")
	ToEmail = os.Getenv("TO_EMAIL")
}

func TestSendsEmail(t *testing.T) {
	err := SendEmail(FromEmail, ToEmail)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}
