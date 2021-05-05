package main

import (
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestPollSqsMessage(t *testing.T) {
	defer gock.Off()

	gock.New("https://sqs.eu-west-2.amazonaws.com").
		Post("/").
		Reply(http.StatusOK).
		BodyString(`{}`)

	PollSqsMessage()
}

func TestRetryPollSqsMessage(t *testing.T) {
	defer gock.Off()

	gock.New("https://sqs.eu-west-2.amazonaws.com").
		Post("/").
		Reply(http.StatusOK).
		BodyString(`{}`)

	RetryPollSqsMessage()
}
