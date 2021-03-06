package main

import (
	"os"
	"testing"
)

var FromEmail string
var ToEmail string

func TestInit(t *testing.T) {
	FromEmail = os.Getenv("FROM_EMAIL")
	ToEmail = os.Getenv("TO_EMAIL")
	if ToEmail == "" {
		t.Fail()
	}
	if FromEmail == "" {
		t.Fail()
	}
}

func TestMessageDefaults(t *testing.T) {
	msg := &Message{}
	if msg.Html != "" {
		t.Errorf("Html was not false: %s", msg.Html)
	}
	if msg.Track != false {
		t.Errorf("Track was not false: %t", msg.Track)
	}
}

func TestRequiresTo(t *testing.T) {
	service := NewMailService()
	message := Message{From: FromEmail, Subject: "Queued Message", Text: "This tests that messages can be queued"}
	err := service.SendEmail(&message)
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "missing param: to address not set" {
		t.FailNow()
	}
}

func TestSendsEmail(t *testing.T) {
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, Subject: "Queued Message", Text: "This tests that messages can be queued"}
	err := service.SendEmail(&message)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestSendsHtmlEmail(t *testing.T) {
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, Subject: "Queued Message", Text: "This tests that messages can be queued", Html: "This tests that emails can contain html"}
	err := service.SendEmail(&message)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestSendsEmailWithReplyTo(t *testing.T) {
	service := NewMailService()
	message := Message{From: FromEmail, ReplyTo: FromEmail, To: ToEmail, Subject: "Queued Message", Text: "This tests that messages can be queued"}
	err := service.SendEmail(&message)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestFromEmailBlocked(t *testing.T) {
	os.Setenv("BLOCK_LIST_ADDRESSES", `\.co\.uk$`)
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, Subject: "Blocked Message", Text: "This tests that messages can be blocked"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLOCK_LIST_ADDRESSES")
	if err == nil {
		t.FailNow()
	}
	if "blocked email" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestReplyToEmailBlocked(t *testing.T) {
	os.Setenv("BLOCK_LIST_ADDRESSES", `mailinator\.com$`)
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, ReplyTo: "test@mailinator.com", Subject: "Blocked Message", Text: "This tests that messages can be blocked"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLOCK_LIST_ADDRESSES")
	if err == nil {
		t.FailNow()
	}
	if "blocked email" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestSenderEmailIsBlocked(t *testing.T) {
	os.Setenv("BLOCK_LIST_ADDRESSES", `tester@mailinator\.com$`)
	service := NewMailService()
	message := Message{From: FromEmail, To: "tester@mailinator.com", ReplyTo: "test@mailinator.com", Subject: "Blocked Message", Text: "This tests that messages can be blocked"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLOCK_LIST_ADDRESSES")
	if err == nil {
		t.FailNow()
	}
	if "blocked sender email" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestContentEmailBlocked(t *testing.T) {
	os.Setenv("BLOCK_LIST_CONTENT", `blocked|test`)
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, ReplyTo: "test@mailinator.com", Subject: "This message will be blocked", Text: "This tests that messages can be blocked when keywords are being filtered"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLOCK_LIST_CONTENT")
	if err == nil {
		t.FailNow()
	}
	if "blocked phrase" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestContentEmailBlockedIgnoreTrailingPipe(t *testing.T) {
	blockList := []string{
		"this string will not get blocked",
	}

	os.Setenv("BLOCK_LIST_CONTENT", `|a|`)
	service := NewMailService()
	for _, element := range blockList {
		message := Message{From: FromEmail, To: ToEmail, ReplyTo: "test@mailinator.com", Subject: "This is a subject", Text: element}
		err := service.SendEmail(&message)
		os.Unsetenv("BLOCK_LIST_CONTENT")
		if err != nil {
			t.FailNow()
		}
	}
}

func TestContentEmailBlockedIgnoresCase(t *testing.T) {
	blockList := []string{
		"lowercase blocked string",
		"Title Case Blocked String",
		"UPPERCASE BLOCKED STRING",
		"RaNdOm CaSe BlOcKeD sTrInG",
		"THIS IS ANOTHER STRING",
	}

	os.Setenv("BLOCK_LIST_CONTENT", `blocked|other`)
	service := NewMailService()

	for _, element := range blockList {
		message := Message{From: FromEmail, To: ToEmail, ReplyTo: "test@mailinator.com", Subject: "This is a subject", Text: element}
		err := service.SendEmail(&message)
		os.Unsetenv("BLOCK_LIST_CONTENT")
		if err == nil {
			t.FailNow()
		}
		if "blocked phrase" != err.Error() {
			t.Errorf("Error: %s", err.Error())
			t.FailNow()
		}
	}
}

func TestContentEmailBlockedRegex(t *testing.T) {
	blockList := []string{
		"Two Words",
		"TwoTogether",
		"&lt;a href=https",
	}

	os.Setenv("BLOCK_LIST_CONTENT", `a href=|two ?(words|together)|single`)
	service := NewMailService()

	for _, element := range blockList {
		message := Message{From: FromEmail, To: ToEmail, ReplyTo: "test@mailinator.com", Subject: "This is a subject", Text: element}
		err := service.SendEmail(&message)
		os.Unsetenv("BLOCK_LIST_CONTENT")
		if err == nil {
			t.FailNow()
		}
		if "blocked phrase" != err.Error() {
			t.Errorf("Error: %s", err.Error())
			t.FailNow()
		}
	}
}

func TestContentEmailEmptyBlockList(t *testing.T) {
	os.Setenv("BLOCK_LIST_CONTENT", ``)
	service := NewMailService()
	message := Message{From: FromEmail, ReplyTo: FromEmail, To: ToEmail, Subject: "Queued Message", Text: "This tests that messages can be queued"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLOCK_LIST_CONTENT")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestContentEmailBlockedUrl(t *testing.T) {
	os.Setenv("BLOCK_LIST_CONTENT", `http[^ ]+http`)
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, ReplyTo: "test@mailinator.com", Subject: "Blocked Message", Text: "This tests that messages can be http://google.com/somewhere?http://blocked.com/address when keywords are being filtered"}
	err := service.SendEmail(&message)
	os.Unsetenv("BLOCK_LIST_CONTENT")
	if err == nil {
		t.FailNow()
	}
	if "blocked phrase" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestSendsEmailIpBlocked(t *testing.T) {
	service := NewMailService()
	message := Message{From: FromEmail, To: ToEmail, Subject: "Queued Message", Text: "This tests that messages can be queued", SenderIp: "185.104.184.126"}
	err := service.SendEmail(&message)
	if err == nil {
		t.FailNow()
	}
	if "blocked ip" != err.Error() {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}

func TestQueuesEmail(t *testing.T) {
	service := NewMailService()

	if nil == service.Messages {
		t.Error("Message queue has not been set")
		t.FailNow()
	}

	message := Message{From: FromEmail, To: ToEmail, Subject: "Queued Message", Text: "This tests that messages can be queued"}
	service.QueueEmail(&message)

	response := <-service.Messages
	if response.Text != message.Text {
		t.FailNow()
	}
}

func TestSendTemplatedEmail(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, Subject: "Queued Message", Template: "test"}
	err := service.SendTemplatedEmail(&message)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
		t.FailNow()
	}
}

func TestSendTemplatedEmailRequiresTo(t *testing.T) {
	service := NewMailService()
	message := Message{Subject: "Queued Message", Template: "test"}
	err := service.SendTemplatedEmail(&message)
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "missing param: to address not set" {
		t.FailNow()
	}
}

func TestSendTemplatedEmailWithBlockedContent(t *testing.T) {
	blockList := []string{
		"Two Words",
		"TwoTogether",
		"&lt;a href=https",
	}

	os.Setenv("BLOCK_LIST_CONTENT", `a href=|two ?(words|together)|single`)
	service := NewMailService()

	for _, element := range blockList {
		params := map[string]interface{}{"content": element}
		message := Message{To: ToEmail, Subject: "Queued Message", Template: "test", Params: params}
		err := service.SendTemplatedEmail(&message)
		os.Unsetenv("BLOCK_LIST_CONTENT")
		if err == nil {
			t.Errorf("Expecting a blocked phrase error")
			t.FailNow()
		}
		if "blocked phrase" != err.Error() {
			t.Errorf("Error: %s", err.Error())
			t.FailNow()
		}
	}
}

func TestQueueTemplatedEmail(t *testing.T) {
	service := NewMailService()

	message := Message{From: FromEmail, To: ToEmail, Subject: "Queued Message", Template: "test"}
	service.QueueEmail(&message)

	response := <-service.Messages
	if response.Template != message.Template {
		t.FailNow()
	}
}

func TestAssembleTemplateWithInvalidTemplate(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, Template: "nonexistant"}
	err := service.AssembleTemplate(&message)
	if err == nil {
		t.FailNow()
	}
}

func TestAssembleTemplate(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, Template: "test"}
	service.AssembleTemplate(&message)
	if message.From == "" {
		t.Errorf("Error: from for email was not assembled\n")
		t.FailNow()
	}
	if message.Subject == "" {
		t.Errorf("Error: subject for email was not assembled\n")
		t.FailNow()
	}
	if message.Text == "" {
		t.Errorf("Error: text for email was not assembled\n")
		t.FailNow()
	}
	if message.Html == "" {
		t.Errorf("Error: html for email was not assembled\n")
		t.FailNow()
	}
}

func TestAssembleTemplateWithReplyTo(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, ReplyTo: FromEmail, Template: "test"}
	service.AssembleTemplate(&message)
	if message.ReplyTo == "" {
		t.Errorf("Error: reply to for email was not assembled\n")
		t.FailNow()
	}
}

func TestAssembleTemplateDoesNotChangeSubject(t *testing.T) {
	service := NewMailService()
	message := Message{To: ToEmail, Subject: "Queued Message", Template: "test"}
	service.AssembleTemplate(&message)
	if message.Subject != "Queued Message" {
		t.Errorf("Error: subject for email was not assembled\n")
		t.FailNow()
	}
}
