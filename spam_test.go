package main

import (
	"os"
	"testing"
)

func TestCheckBlackList(t *testing.T) {
	FromEmail = os.Getenv("FROM_EMAIL")
	ToEmail = os.Getenv("TO_EMAIL")
	m := &Message{
		To:       ToEmail,
		From:     FromEmail,
		Subject:  "Hello world",
		Text:     "Hello world",
		SenderIp: "127.0.0.1",
	}
	CheckBlackList(m.SenderIp)
}

func TestCheckAkismet(t *testing.T) {
	FromEmail = os.Getenv("FROM_EMAIL")
	msg := &Message{
		Text:      "This is a test",
		Website:   "http://www.thomaslorentsen.co.uk",
		UserAgent: "golang",
		FromName:  "Cuthbert Rumbold",
		From:      FromEmail,
		SenderIp:  "127.0.0.1",
	}
	result := CheckAkismet(msg)
	if result == true {
		t.FailNow()
	}
}
