package templates

import (
	"bytes"
	"text/template"
)

type Template struct {
	Name    string
	Content *template.Template
}

func (t *Template) Generate(u *Unit) (string, error) {
	var unitTemplate bytes.Buffer

	err := t.Content.Execute(&unitTemplate, u)
	if err != nil {
		return "", err
	}

	return unitTemplate.String(), nil
}
