package main

import "testing"

func TestSendsEmail(t *testing.T) {
	err := SendEmail()
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
}
