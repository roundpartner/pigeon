package main

import "testing"
import (
	"net/http"
	"net/http/httptest"
	"strings"
)

func TestMailService_SendEmail(t *testing.T) {
	body := strings.NewReader("{\"to\":\"receipient@mailinator.com\",\"from\":\"sender@mailinator.com\",\"subject\":\"Cool Subject\",\"text\":\"Interesting Message\"}")
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/email", body)

	rs := NewRestServer()
	rs.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("Service did not return ok status")
		t.FailNow()
	}
}

func TestMailService_SendEmailFailsWithNoToAddress(t *testing.T) {
	body := strings.NewReader("{\"from\":\"sender@mailinator.com\",\"subject\":\"Cool Subject\",\"text\":\"Interesting Message\"}")
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/email", body)

	rs := NewRestServer()
	rs.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Service did not return status bad request")
		t.FailNow()
	}
}
