package codegen

import (
	"bytes"
	"fmt"
	"go/types"
	"strings"
	"text/template"

	"github.com/sunboyy/repogen/internal/code"
)

const structTemplate = `
type {{.Name}} struct {
{{.GenFields}}
}
`

// StructBuilder is an implementer of a struct.
type StructBuilder struct {
	Pkg    *types.Package
	Name   string
	Fields []code.StructField
}

// Impl writes struct declatation code to the buffer.
func (sb StructBuilder) Impl(buffer *bytes.Buffer) error {
	tmpl, err := template.New("struct").Parse(structTemplate)
	if err != nil {
		return err
	}

	// writing to a buffer should not cause errors.
	_ = tmpl.Execute(buffer, sb)

	return nil
}

func (sb StructBuilder) GenFields() string {
	var fieldLines []string
	for _, field := range sb.Fields {
		fieldLine := fmt.Sprintf("\t%s %s", field.Var.Name(), TypeToString(sb.Pkg, field.Var.Type()))
		if len(field.Tag) > 0 {
			fieldLine += fmt.Sprintf(" `%s`", string(field.Tag))
		}
		fieldLines = append(fieldLines, fieldLine)
	}
	return strings.Join(fieldLines, "\n")
}
