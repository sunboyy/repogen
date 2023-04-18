package codegen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/sunboyy/repogen/internal/code"
)

const functionTemplate = `
func {{.Name}}({{.GenParams}}){{.GenReturns}} {
{{.Body.Code}}
}
`

// FunctionBuilder is an implementer of a function.
type FunctionBuilder struct {
	Name    string
	Params  []code.Param
	Returns []code.Type
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
	return generateParams(fb.Params)
}

func (fb FunctionBuilder) GenReturns() string {
	return generateReturns(fb.Returns)
}

func generateParams(params []code.Param) string {
	var paramList []string
	for _, param := range params {
		paramList = append(
			paramList,
			fmt.Sprintf("%s %s", param.Name, param.Type.Code()),
		)
	}
	return strings.Join(paramList, ", ")
}

func generateReturns(returns []code.Type) string {
	if len(returns) == 0 {
		return ""
	}

	if len(returns) == 1 {
		return " " + returns[0].Code()
	}

	var returnList []string
	for _, ret := range returns {
		returnList = append(returnList, ret.Code())
	}

	return fmt.Sprintf(" (%s)", strings.Join(returnList, ", "))
}
