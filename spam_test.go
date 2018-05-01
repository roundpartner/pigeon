package main

import (
	"os"
	"testing"
)

func TestCheckSpamAssassin(t *testing.T) {
	FromEmail = os.Getenv("FROM_EMAIL")
	ToEmail = os.Getenv("TO_EMAIL")
	m := &Message{
		To:      ToEmail,
		From:    FromEmail,
		Subject: "Hello world",
		Text:    "Hello world",
	}
	CheckSpamAssassin(m)
}
