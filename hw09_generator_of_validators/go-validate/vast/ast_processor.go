package vast

import (
	"go/ast"
	"strings"
)

const ValidateTagWord = "validate"

type Validation map[string][]ValidationField

type ValidationField struct {
	Name     string
	Type     string
	BaseType string
	Rule     ValidationRule
}

type ValidationRule struct {
	Name  string
	Value string
}

func Parse(node *ast.File) (Validation, error) { //nolint:gocognit
	var validation = make(Validation)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, specs := range genDecl.Specs {
			typeSpecs, ok := specs.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpecs.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range structType.Fields.List {
				if field.Tag == nil {
					continue
				}

				tagValue := parseTag(field.Tag.Value)
				if tagValue == "" {
					continue
				}

				fieldType := getFieldType(field)
				if fieldType == "" {
					continue
				}
				baseType, ok := getBaseType(node, getFieldType(field))
				if !ok || baseType == "" {
					continue
				}

				// all checks are passed
				addToValidate(validation, typeSpecs.Name.String(), field.Names[0].Name, fieldType, baseType, tagValue)
			}
		}
	}

	return validation, nil
}

func addToValidate(m Validation, structName, fieldName, fieldType, fieldBaseType, tagValue string) {
	rules := strings.Split(tagValue, "|")
	for _, v := range rules {
		rule := strings.Split(v, ":")
		m[structName] = append(m[structName], ValidationField{
			Name:     fieldName,
			Type:     fieldType,
			BaseType: fieldBaseType,
			Rule: ValidationRule{
				Name:  rule[0],
				Value: rule[1],
			},
		})
	}
}

func parseTag(tag string) string {
	if tag == "" {
		return ""
	}
	s := strings.Split(strings.Trim(tag, "`"), " ")
	for _, v := range s {
		if strings.HasPrefix(v, ValidateTagWord) {
			return strings.Trim(strings.TrimPrefix(v, ValidateTagWord+":"), `"`)
		}
	}
	return ""
}

func getFieldType(f *ast.Field) string {
	if f == nil {
		return ""
	}
	switch t := f.Type.(type) {
	case *ast.Ident:
		return t.String()

	case *ast.ArrayType:
		tt, ok := t.Elt.(*ast.Ident)
		if !ok {
			return ""
		}
		return "[]" + tt.String()
	}
	return ""
}

// нужно для получения базового типа.
func getBaseType(node ast.Node, typeName string) (string, bool) {
	if typeName == "string" || typeName == "int" {
		return typeName, true
	}
	if strings.HasPrefix(typeName, "[]") {
		if typeName == "[]string" || typeName == "[]int" {
			return typeName, true
		}
	}

	switch strings.HasPrefix(typeName, "[]") {
	case true:
		t := inspectBaseType(node, strings.TrimPrefix(typeName, "[]"))
		if t == "" {
			return "", false
		}
		if strings.HasPrefix(t, "[]") {
			return "", false
		}
		return "[]" + t, true
	case false:
		t := inspectBaseType(node, typeName)
		if t == "" {
			return "", false
		}
		return t, true
	}
	return "", false
}

// нужно для получения базового типа.
func inspectBaseType(node ast.Node, typeName string) string {
	var stop bool
	var result string

	ast.Inspect(node, func(n ast.Node) bool {
		if stop {
			return false
		}
		tp, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if tp.Name.String() == typeName {
			switch t := tp.Type.(type) {
			case *ast.Ident:
				if t.String() == "string" {
					result = t.String()
					stop = true
					return false
				}
				if t.String() == "int" {
					result = t.String()
					stop = true
					return false
				}
			case *ast.ArrayType:
				tt, ok := t.Elt.(*ast.Ident)
				if !ok {
					stop = true
					return false
				}
				result = "[]" + tt.String()
				stop = true
				return false
			}
		}
		return true
	})

	return result
}
