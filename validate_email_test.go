package main

import "testing"

func TestValidateEmail(t *testing.T) {
	result := ValidateEmail("tester@mailinator.com")
	if false == result {
		t.FailNow()
	}
}

func TestValidateEmailInvalidAddress(t *testing.T) {
	result := ValidateEmail("mailinator.com")
	if true == result {
		t.FailNow()
	}
}
