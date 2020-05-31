package main

import (
	"io"
	"strings"
	"text/template"

	"github.com/fdully/hw_otus/hw09_generator_of_validators/go-validate/vast"
)

type TemplateData struct {
	PackageName string
	Validation  vast.Validation
}

var funcMap = template.FuncMap{"addQuotes": splitStringAndAddQuotes}

var (
	//funcMap = templates.FuncMap{"getValidationType": getValidationType, "getLenValidationValue": getLenValidationValue}

	validateTmpl = template.Must(template.New("validate").Funcs(funcMap).Parse(`
         // Code generated by cool go-validate tool; DO NOT EDIT.
         package {{ .PackageName }}
        
         import (
        	"unicode/utf8" 
            "regexp"
            "errors"
         )

         type ValidationError struct {
	       Field string
	       Err   error
         }
       
        {{ range $name, $fields := .Validation }}
        func (_this {{ $name }}) Validate() ([]ValidationError, error) {
		var val []ValidationError 

        {{- range $fields }}
          {{- if eq .BaseType "int" }}
            {{ template "int" . }}
          {{- end }}

          {{- if eq .BaseType "[]int" }}
          for _, v := range _this.{{ .Name }} {
            {{ template "[]int" . }}
          }
          {{- end }}

        {{- if eq .BaseType "string" }}
          {{ template "string" . }}
        {{- end }}

        {{- if eq .BaseType "[]string" }}
          for _, v := range _this.{{ .Name }} {
            {{ template "[]string" . }}
          }
        {{- end }}

        {{- end }}
           return val, nil
        }
        {{- end }}
    `))

	_ = template.Must(validateTmpl.New("int").Parse(`
        {{- if eq .Rule.Name "min" }}
          if {{ .Rule.Value }} > _this.{{ .Name }} {
            val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
        {{- end }}
        {{- if eq .Rule.Name "max" }}
          if {{ .Rule.Value }} < _this.{{ .Name }} {
            val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
        {{- end }}
        {{- if eq .Rule.Name "in" }}
          { var isIn bool
          for _, v := range []{{ .Type }}{ {{ .Rule.Value }} } {
            if v == _this.{{ .Name }} {
              isIn = true
            }
          }
          if !isIn {
              val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
          }
        {{- end }}
     `))

	_ = template.Must(validateTmpl.New("[]int").Parse(`
        {{- if eq .Rule.Name "max" }}
          if v > {{ .Rule.Value }} {
            val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
        {{- end }}
        {{- if eq .Rule.Name "min" }}
          if v < {{ .Rule.Value }} {
            val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
        {{- end }}
     `))

	_ = template.Must(validateTmpl.New("string").Parse(`
        {{- if eq .Rule.Name "len" }}
          if utf8.RuneCountInString(_this.{{ .Name }}) != {{ .Rule.Value }} {
            val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
        {{- end }}
        {{- if eq .Rule.Name "regexp" }}
          re := regexp.MustCompile("{{ .Rule.Value }}")
          if !re.MatchString(_this.{{ .Name }}) {
            val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
        {{- end }}
        {{- if eq .Rule.Name "in" }}
          { var isIn bool
          for _, v := range []{{ .Type }}{ {{ addQuotes .Rule.Value }} } {
            if v == _this.{{ .Name }} {
              isIn = true
            }
          }
          if !isIn {
              val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
          }
        {{- end }}
     `))

	_ = template.Must(validateTmpl.New("[]string").Parse(`
        {{- if eq .Rule.Name "len" }}
          if utf8.RuneCountInString(v) != {{ .Rule.Value }} {
            val = append(val, ValidationError{Field: "{{ .Name }}", Err: errors.New("value must be {{ .Rule.Name }} {{ .Rule.Value }}")})
          }
        {{- end }}
     `))
)

func ExecuteTemplate(tmpl *template.Template, wr io.Writer, data TemplateData) error {
	return tmpl.Execute(wr, data)
}

func splitStringAndAddQuotes(str string) string {
	if str == "" {
		return ""
	}
	var s = make([]string, 0, 1)
	for _, v := range strings.Split(str, ",") {
		s = append(s, `"`+v+`"`)
	}
	return strings.Join(s, ",")
}