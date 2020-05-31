package main

import (
	"github.com/fdully/hw_otus/hw09_generator_of_validators/go-validate/vast"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestExecuteTemplate(t *testing.T) {

	r := vast.Validation{"App": []vast.ValidationField{{Name: "Version", Type: "string", Rule: vast.ValidationRule{Name: "len", Value: "13"}},
		{Name: "Code", Type: "int", Rule: vast.ValidationRule{Name: "in", Value: "200,300,400"}},
		{Name: "Num", Type: "int", Rule: vast.ValidationRule{Name: "max", Value: "13"}},
		{Name: "Num", Type: "int", Rule: vast.ValidationRule{Name: "min", Value: "3"}},
		{Name: "Names", Type: "[]string", Rule: vast.ValidationRule{Name: "len", Value: "5"}}}}

	tmplData := TemplateData{
		PackageName: "models",
		Validation:  r,
	}

	require.NoError(t, ExecuteTemplate(validateTmpl, ioutil.Discard, tmplData))
}

func TestInRuleFuncMap(t *testing.T) {
	require.Equal(t, "\"go\"", splitStringAndAddQuotes("go"))
	require.Equal(t, "\"go\",\"go\"", splitStringAndAddQuotes("go,go"))
	require.Equal(t, "", splitStringAndAddQuotes(""))
}
