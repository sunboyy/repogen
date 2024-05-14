package codegen

import (
	"bytes"
	"fmt"
	"go/types"
	"strings"
	"text/template"
)

const functionTemplate = `
func {{.Name}}({{.GenParams}}){{.GenReturns}} {
{{.Body.Code}}
}
`

// FunctionBuilder is an implementer of a function.
type FunctionBuilder struct {
	Pkg     *types.Package
	Name    string
	Params  *types.Tuple
	Returns []types.Type
	Body    FunctionBody
}

// Impl writes function declatation code to the buffer.
func (fb FunctionBuilder) Impl(buffer *bytes.Buffer) error {
	tmpl, err := template.New("function").Parse(functionTemplate)
	if err != nil {
		return err
	}

	// writing to a buffer should not cause errors.
	_ = tmpl.Execute(buffer, fb)

	return nil
}

func (fb FunctionBuilder) GenParams() string {
	return generateParams(fb.Pkg, fb.Params)
}

func (fb FunctionBuilder) GenReturns() string {
	return generateReturns(fb.Pkg, fb.Returns)
}

func generateParams(pkg *types.Package, params *types.Tuple) string {
	var paramList []string
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)

		paramList = append(
			paramList,
			fmt.Sprintf("%s %s", param.Name(), typeToString(pkg, param.Type())),
		)
	}
	return strings.Join(paramList, ", ")
}

func typeToString(pkg *types.Package, t types.Type) string {
	switch t := t.(type) {
	case *types.Pointer:
		return fmt.Sprintf("*%s", typeToString(pkg, t.Elem()))

	case *types.Slice:
		return fmt.Sprintf("[]%s", typeToString(pkg, t.Elem()))

	case *types.Named:
		if t.Obj().Pkg() == nil || t.Obj().Pkg().Path() == pkg.Path() {
			return t.Obj().Name()
		}
		return fmt.Sprintf("%s.%s", t.Obj().Pkg().Name(), t.Obj().Name())

	default:
		return t.String()
	}
}

func generateReturns(pkg *types.Package, returns []types.Type) string {
	if len(returns) == 0 {
		return ""
	}

	if len(returns) == 1 {
		return " " + typeToString(pkg, returns[0])
	}

	var returnList []string
	for _, ret := range returns {
		returnList = append(returnList, typeToString(pkg, ret))
	}

	return fmt.Sprintf(" (%s)", strings.Join(returnList, ", "))
}
