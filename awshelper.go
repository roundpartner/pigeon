package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

var awssession = &session.Session{}

func GetAWSSession() *session.Session {
	region := os.Getenv("AWS_REGION")
	if _, exists := os.LookupEnv("AWS_ACCESS_KEY_ID"); exists {
		session, err := session.NewSession(
			&aws.Config{
				Region:      aws.String(region),
				Credentials: credentials.NewEnvCredentials(),
			})
		if err != nil {
			log.Printf("AWS Error: %s", err.Error())
			return nil
		}
		return session
	}
	session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
		})
	if err != nil {
		log.Printf("AWS Error: %s", err.Error())
		return nil
	}
	return session
}

func GetQueueName() (string, error) {
	queue, exists := os.LookupEnv("AWS_SQS_QUEUE")
	if exists == false {
		log.Printf("[ERROR] [%s] %s", ServiceName, "Queue not set")
		return "", errors.New("queue not set")
	}
	if queue == "" {
		log.Printf("[ERROR] [%s] %s", ServiceName, "Queue not set")
		return "", errors.New("queue not set")
	}
	return queue, nil
}
