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

func TestMailService_SendEmailWithReport(t *testing.T) {
	body := strings.NewReader("{\"to\":\"receipient@mailinator.com\",\"from\":\"sender@mailinator.com\",\"subject\":\"Cool Subject\",\"text\":\"Interesting Message\",\"report\":true}")
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

func TestMailService_ViewTemplatedEmail(t *testing.T) {
	body := strings.NewReader("{\"template\":\"test\",\"params\": {\"name\": \"Cuthbert\",\"colour\": \"Purple\"}}")
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/template", body)

	rs := NewRestServer()
	rs.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Service did not return ok status")
		t.FailNow()
	}
	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fatalf("Service did not return json header")
		t.FailNow()
	}
	if "" == rr.Body.String() {
		t.Fatalf("Empty body returned: %s", rr.Body.String())
		t.FailNow()
	}
}

func TestVerifyIpIsBlocked(t *testing.T) {
	body := strings.NewReader("{\"ip\":\"185.104.184.126\"}")
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/verify", body)

	rs := NewRestServer()
	rs.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Service did not return ok status")
		t.FailNow()
	}
	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fatalf("Service did not return json header")
		t.FailNow()
	}
	if `{"ip":"185.104.184.126","blocked":true}` != rr.Body.String() {
		t.Fatalf("Unexpected body returned: %s", rr.Body.String())
		t.FailNow()
	}
}

func TestVerifyIpIsNotBlocked(t *testing.T) {
	body := strings.NewReader("{\"ip\":\"127.0.0.1\"}")
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/verify", body)

	rs := NewRestServer()
	rs.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Service did not return ok status")
		t.FailNow()
	}
	if "application/json; charset=utf-8" != rr.Header().Get("Content-Type") {
		t.Fatalf("Service did not return json header")
		t.FailNow()
	}
	if `{"ip":"127.0.0.1","blocked":false}` != rr.Body.String() {
		t.Fatalf("Unexpected body returned: %s", rr.Body.String())
		t.FailNow()
	}
}
