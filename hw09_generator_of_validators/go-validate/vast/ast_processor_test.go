package vast

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

var testSrcAst = `
package models

type UserRole string

type App struct {
		Version string ` + fmt.Sprint("`json:\"id\" validate:\"len:13\"`") + `
        Code []int ` + fmt.Sprint("`validate:\"in:200,300,400\"`") + `
        Num int ` + fmt.Sprint("`validate:\"max:13|min:3\"`") + `
        Names UserRole ` + fmt.Sprint("`validate:\"in:admin,user\"`") + `
        Struct struct{}
        Map map[string]string ` + fmt.Sprint("`validate:\"len:13\"`") + `
		EmptyTag string
	}
`

func TestGetFieldType(t *testing.T) {

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", testSrcAst, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	fieldAppVersion := node.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
	fieldAppCode := node.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[1]
	fieldAppNum := node.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[2]
	fieldAppNames := node.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[3]
	fieldAppStruct := node.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[4]
	fieldAppMap := node.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[5]

	version, ok := getBaseType(node, getFieldType(fieldAppVersion))
	require.Equal(t, "string", version)
	require.Equal(t, true, ok)

	code, ok := getBaseType(node, getFieldType(fieldAppCode))
	require.Equal(t, "[]int", code)
	require.Equal(t, true, ok)

	num, ok := getBaseType(node, getFieldType(fieldAppNum))
	require.Equal(t, "int", num)
	require.Equal(t, true, ok)

	names, ok := getBaseType(node, getFieldType(fieldAppNames))
	require.Equal(t, "string", names)
	require.Equal(t, true, ok)

	structType, ok := getBaseType(node, getFieldType(fieldAppStruct))
	require.Equal(t, "", structType)
	require.Equal(t, false, ok)

	mapType, ok := getBaseType(node, getFieldType(fieldAppMap))
	require.Equal(t, "", mapType)
	require.Equal(t, false, ok)

}

func TestParseTag(t *testing.T) {
	tag := `json:"id" validate:"min=5"`
	require.Equal(t, "min=5", parseTag(tag))

	tag = `validate:"max=4|len=3"`
	require.Equal(t, "max=4|len=3", parseTag(tag))

	tag = `validate:"in:200,300,400"`
	require.Equal(t, "in:200,300,400", parseTag(tag))

	tag = `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	require.Equal(t, `regexp:^\\w+@\\w+\\.\\w+$`, parseTag(tag))

	tag = `json:"id"`
	require.Equal(t, "", parseTag(tag))

	tag = ""
	require.Equal(t, "", parseTag(tag))
}

func TestParse(t *testing.T) {

	r := Validation{"App": []ValidationField{{Name: "Version", Type: "string", BaseType: "string", Rule: ValidationRule{Name: "len", Value: "13"}},
		{Name: "Code", Type: "[]int", BaseType: "[]int", Rule: ValidationRule{Name: "in", Value: "200,300,400"}},
		{Name: "Num", Type: "int", BaseType: "int", Rule: ValidationRule{Name: "max", Value: "13"}},
		{Name: "Num", Type: "int", BaseType: "int", Rule: ValidationRule{Name: "min", Value: "3"}},
		{Name: "Names", Type: "UserRole", BaseType: "string", Rule: ValidationRule{Name: "in", Value: "admin,user"}}}}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", testSrcAst, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	rr, err := Parse(node)
	require.NoError(t, err)
	require.Equal(t, r, rr)

}
