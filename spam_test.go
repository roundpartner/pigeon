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
