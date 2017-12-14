package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestAssemble(t *testing.T) {
	m := map[string]interface{}{}
	b := bytes.NewBufferString("{\"name\":\"Cuthbert\"}").Bytes()
	json.Unmarshal(b, &m)

	text, err := Assemble("testing", "Hello {{ .name }}", &m)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
	if text != "Hello Cuthbert" {
		t.Errorf("Text did not match: %s", text)
	}
}

func TestAssembleFromMessage(t *testing.T) {
	b := bytes.NewBufferString("{\"params\":{\"name\":\"Cuthbert\", \"food\":\"Cake\"}}")
	decoder := json.NewDecoder(b)
	msg := &Message{}
	err := decoder.Decode(msg)
	if err != nil {
		t.Errorf("Unable to decode: %s", err.Error())
		t.FailNow()
	}

	tpl := &Template{"testing", "Hello {{ .name }} and my favourite food is {{ .food }}"}
	text, err := AssembleTemplate(tpl, msg)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
	if text != "Hello Cuthbert and my favourite food is Cake" {
		t.Errorf("Text did not match: %s", text)
	}
}

func TestImportEmailTemplate(t *testing.T) {
	tm := NewTemplateManager()
	tpl, err := tm.ImportTemplate("test")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.FailNow()
	}
	if tpl.Subject == "" {
		t.Errorf("Subject not set")
		t.Fail()
	}
	if tpl.Text == nil || tpl.Text.Name == "" {
		t.Errorf("Text template not set")
		t.Fail()
	}
	if tpl.Html == nil || tpl.Html.Name == "" {
		t.Errorf("Html template not set")
		t.Fail()
	}
}
