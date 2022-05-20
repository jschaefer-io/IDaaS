package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

type TemplateList struct {
	templates *template.Template
}

func (t *TemplateList) Execute(writer io.Writer, namespace string, name string, data any) error {
	return t.templates.ExecuteTemplate(writer, fmt.Sprintf("%s-%s.gohtml", namespace, name), data)
}

func (t *TemplateList) ExecuteToString(namespace string, name string, data any) (string, error) {
	var buf bytes.Buffer
	err := t.Execute(&buf, namespace, name, data)
	return buf.String(), err
}

func NewTemplateList(dir string) *TemplateList {
	return &TemplateList{
		templates: template.Must(template.ParseGlob(dir + "/*.gohtml")),
	}
}
