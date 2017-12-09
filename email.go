package main

import (
	"gopkg.in/mailgun/mailgun-go.v1"
	"fmt"
	"os"
)

func SendEmail(from string, to string) error {
	domain := os.Getenv("DOMAIN")
	apiKey := os.Getenv("API_KEY")
	publicApiKey := os.Getenv("PUBLIC_API_KEY")

	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	message := mg.NewMessage(
		from,
		"Fancy subject!",
		"Hello from Mailgun Go!",
		to)
	resp, id, err := mg.Send(message)

	if err != nil {
		return err
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return err
}
