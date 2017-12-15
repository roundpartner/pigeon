package main

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"text/template"
)

type Template struct {
	Name string `json:"name" yaml:"name"`
	Text string `json:"text" yaml:"text"`
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
	From    string    `json:"from" yaml:"from"`
	Subject string    `json:"subject" yaml:"subject"`
	Text    *Template `json:"text" yaml:"text"`
	Html    *Template `json:"html" yaml:"html"`
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
	b, err := ioutil.ReadFile(tm.TemplateDirectory + "/" + name + ".yml")
	if err != nil {
		return nil, err
	}
	tpl := &EmailTemplate{}
	err = yaml.Unmarshal(b, tpl)
	return tpl, err
}
