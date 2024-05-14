package codegen

import (
	"bytes"
	"fmt"
	"go/types"
	"text/template"

	"github.com/sunboyy/repogen/internal/code"
)

const methodTemplate = `
func ({{.GenReceiver}}) {{.Name}}({{.GenParams}}){{.GenReturns}} {
{{.Body.Code}}
}
`

// MethodBuilder is an implementer of a method.
type MethodBuilder struct {
	Pkg      *types.Package
	Receiver MethodReceiver
	Name     string
	Params   *types.Tuple
	Returns  []types.Type
	Body     FunctionBody
}

// MethodReceiver describes a specification of a method receiver.
type MethodReceiver struct {
	Name    string
	Type    code.SimpleType
	Pointer bool
}

// Impl writes method declatation code to the buffer.
func (mb MethodBuilder) Impl(buffer *bytes.Buffer) error {
	tmpl, err := template.New("function").Parse(methodTemplate)
	if err != nil {
		return err
	}

	// writing to a buffer should not cause errors.
	_ = tmpl.Execute(buffer, mb)

	return nil
}

func (mb MethodBuilder) GenReceiver() string {
	if mb.Receiver.Name == "" {
		return mb.generateReceiverType()
	}
	return fmt.Sprintf("%s %s", mb.Receiver.Name, mb.generateReceiverType())
}

func (mb MethodBuilder) generateReceiverType() string {
	if !mb.Receiver.Pointer {
		return mb.Receiver.Type.Code()
	}
	return code.PointerType{ContainedType: mb.Receiver.Type}.Code()
}

func (mb MethodBuilder) GenParams() string {
	return generateParams(mb.Pkg, mb.Params)
}

func (mb MethodBuilder) GenReturns() string {
	return generateReturns(mb.Pkg, mb.Returns)
}
