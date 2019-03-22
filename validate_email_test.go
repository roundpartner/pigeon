package main

import "testing"

func TestValidateEmail(t *testing.T) {
	result := ValidateEmail("tom@thomaslorentsen.co.uk")
	if false == result {
		t.FailNow()
	}
}

func TestValidateEmailInvalidAddress(t *testing.T) {
	result := ValidateEmail("thomaslorentsen.co.uk")
	if true == result {
		t.FailNow()
	}
}
