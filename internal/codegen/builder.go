package codegen

import (
	"bytes"
	"text/template"

	"github.com/sunboyy/repogen/internal/code"
	"golang.org/x/tools/imports"
)

type Builder struct {
	// Program defines generator program name in the generated file.
	Program string

	// PackageName defines the package name of the generated file.
	PackageName string

	// Imports defines necessary imports to reduce ambiguity when generating
	// formatting the raw-generated code.
	Imports [][]code.Import

	implementers []Implementer
}

// Implementer is an interface that wraps the basic Impl method for code
// generation.
type Implementer interface {
	Impl(buffer *bytes.Buffer) error
}

// NewBuilder is a constructor of Builder struct.
func NewBuilder(program string, packageName string, imports [][]code.Import) *Builder {
	return &Builder{
		Program:     program,
		PackageName: packageName,
		Imports:     imports,
	}
}

// AddImplementer appends a new implemeneter to the implementer list.
func (b *Builder) AddImplementer(implementer Implementer) {
	b.implementers = append(b.implementers, implementer)
}

// Build generates code from the previously provided specifications.
func (b Builder) Build() (string, error) {
	buffer := new(bytes.Buffer)

	if err := b.buildBase(buffer); err != nil {
		return "", err
	}

	for _, impl := range b.implementers {
		if err := impl.Impl(buffer); err != nil {
			return "", err
		}
	}

	formattedCode, err := imports.Process("", buffer.Bytes(), nil)
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func (b Builder) buildBase(buffer *bytes.Buffer) error {
	tmpl, err := template.New("file_base").Parse(baseTemplate)
	if err != nil {
		return err
	}

	tmplData := baseTemplateData{
		Program:     b.Program,
		PackageName: b.PackageName,
		Imports:     b.Imports,
	}

	// writing to a buffer should not cause errors.
	_ = tmpl.Execute(buffer, tmplData)

	return nil
}
