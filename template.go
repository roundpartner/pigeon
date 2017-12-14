package main

import (
	"bytes"
	"text/template"
)

type Template struct {
	Name string
	Text string
}

func AssembleTemplate(tpl *Template, msg *Message) (string, error) {
	return Assemble(tpl.Name, tpl.Text, msg.Params)
}

func Assemble(name string, text string, data interface{}) (string, error) {
	t, err := template.New(name).Parse(text)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
