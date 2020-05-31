package main

import (
	"bytes"
	"go/ast"
	"go/format"

	"github.com/fdully/hw_otus/hw09_generator_of_validators/go-validate/vast"
)

func GenerateValidationData(node *ast.File) ([]byte, error) {
	validationData, err := vast.Parse(node)
	if err != nil {
		return nil, err
	}

	var tmplData TemplateData
	tmplData.PackageName = node.Name.Name
	tmplData.Validation = validationData

	buf := new(bytes.Buffer)
	err = ExecuteTemplate(validateTmpl, buf, tmplData)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
