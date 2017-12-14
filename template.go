package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"
)

type Template struct {
	Name string `json:"name"`
	Text string `json:"text"`
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

type EmailTemplate struct {
	From    string    `json:"from"`
	Subject string    `json:"subject"`
	Text    *Template `json:"text"`
	Html    *Template `json:"html"`
}

func NewTemplateManager() *TemplateManager {
	templateDirectory := "templates"
	if os.Getenv("TEMPLATES") != "" {
		templateDirectory = os.Getenv("TEMPLATES")
	}
	return &TemplateManager{TemplateDirectory: templateDirectory}
}

type TemplateManager struct {
	TemplateDirectory string
}

func (tm *TemplateManager) ImportTemplate(name string) (*EmailTemplate, error) {
	b, err := ioutil.ReadFile(tm.TemplateDirectory + "/" + name + ".json")
	if err != nil {
		return nil, err
	}
	tpl := &EmailTemplate{}
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	err = decoder.Decode(tpl)
	return tpl, err
}
