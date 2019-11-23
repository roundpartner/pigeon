package main

import "testing"

func TestGetSession(t *testing.T) {
	session := GetAWSSession()
	if session == nil {
		t.Errorf("AWS Session returned nil")
	}
}

func TestGetQueueName(t *testing.T) {
	_, err := GetQueueName()
	if err != nil {
		t.Errorf("Unexpected Error: %s", err.Error())
	}
}
