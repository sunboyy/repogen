package codegen

import (
	"fmt"
	"strings"

	"github.com/sunboyy/repogen/internal/code"
)

type FunctionBody []Statement

func (b FunctionBody) Code() string {
	var lines []string
	for _, statement := range b {
		stmtLines := statement.CodeLines()
		for _, line := range stmtLines {
			lines = append(lines, fmt.Sprintf("\t%s", line))
		}
	}
	return strings.Join(lines, "\n")
}

type Statement interface {
	CodeLines() []string
}

type RawStatement string

func (stmt RawStatement) CodeLines() []string {
	return []string{string(stmt)}
}

type Identifier string

func (id Identifier) CodeLines() []string {
	return []string{string(id)}
}

type DeclStatement struct {
	Name string
	Type code.Type
}

func (stmt DeclStatement) CodeLines() []string {
	return []string{fmt.Sprintf("var %s %s", stmt.Name, stmt.Type.Code())}
}

type DeclAssignStatement struct {
	Vars   []string
	Values StatementList
}

func (stmt DeclAssignStatement) CodeLines() []string {
	vars := strings.Join(stmt.Vars, ", ")
	lines := stmt.Values.CodeLines()
	lines[0] = fmt.Sprintf("%s := %s", vars, lines[0])
	return lines
}

type AssignStatement struct {
	Vars   []string
	Values StatementList
}

func (stmt AssignStatement) CodeLines() []string {
	vars := strings.Join(stmt.Vars, ", ")
	lines := stmt.Values.CodeLines()
	lines[0] = fmt.Sprintf("%s = %s", vars, lines[0])
	return lines
}

type StatementList []Statement

func (l StatementList) CodeLines() []string {
	if len(l) == 0 {
		return []string{""}
	}
	return concatenateStatements(", ", []Statement(l))
}

type ReturnStatement StatementList

func (stmt ReturnStatement) CodeLines() []string {
	lines := StatementList(stmt).CodeLines()
	lines[0] = fmt.Sprintf("return %s", lines[0])
	return lines
}

type ChainStatement []Statement

func (stmt ChainStatement) CodeLines() []string {
	return concatenateStatements(".", []Statement(stmt))
}

type CallStatement struct {
	FuncName string
	Params   StatementList
}

func (stmt CallStatement) CodeLines() []string {
	lines := stmt.Params.CodeLines()
	lines[0] = fmt.Sprintf("%s(%s", stmt.FuncName, lines[0])
	lines[len(lines)-1] += ")"
	return lines
}

type SliceStatement struct {
	Type   code.Type
	Values []Statement
}

func (stmt SliceStatement) CodeLines() []string {
	lines := []string{stmt.Type.Code() + "{"}
	for _, value := range stmt.Values {
		stmtLines := value.CodeLines()
		stmtLines[len(stmtLines)-1] += ","
		for _, line := range stmtLines {
			lines = append(lines, "\t"+line)
		}
	}
	lines = append(lines, "}")
	return lines
}

type MapStatement struct {
	Type  string
	Pairs []MapPair
}

func (stmt MapStatement) CodeLines() []string {
	return generateCollectionCodeLines(stmt.Type, stmt.Pairs)
}

type MapPair struct {
	Key   string
	Value Statement
}

func (p MapPair) ItemCodeLines() []string {
	lines := p.Value.CodeLines()
	lines[0] = fmt.Sprintf(`"%s": %s`, p.Key, lines[0])
	return lines
}

type StructStatement struct {
	Type  string
	Pairs []StructFieldPair
}

func (stmt StructStatement) CodeLines() []string {
	return generateCollectionCodeLines(stmt.Type, stmt.Pairs)
}

type StructFieldPair struct {
	Key   string
	Value Statement
}

func (p StructFieldPair) ItemCodeLines() []string {
	lines := p.Value.CodeLines()
	lines[0] = fmt.Sprintf(`%s: %s`, p.Key, lines[0])
	return lines
}

type collectionItem interface {
	ItemCodeLines() []string
}

func generateCollectionCodeLines[T collectionItem](typ string, pairs []T) []string {
	lines := []string{fmt.Sprintf("%s{", typ)}
	for _, pair := range pairs {
		pairLines := pair.ItemCodeLines()
		pairLines[len(pairLines)-1] += ","
		for _, line := range pairLines {
			lines = append(lines, fmt.Sprintf("\t%s", line))
		}
	}
	lines = append(lines, "}")
	return lines
}

type RawBlock struct {
	Header     []string
	Statements []Statement
}

func (b RawBlock) CodeLines() []string {
	lines := make([]string, len(b.Header))
	copy(lines, b.Header)
	lines[len(lines)-1] += " {"
	for _, stmt := range b.Statements {
		stmtLines := stmt.CodeLines()
		for _, line := range stmtLines {
			lines = append(lines, fmt.Sprintf("\t%s", line))
		}
	}
	lines = append(lines, "}")
	return lines
}

type IfBlock struct {
	Condition  []Statement
	Statements []Statement
}

func (b IfBlock) CodeLines() []string {
	conditionCode := concatenateStatements("; ", b.Condition)
	conditionCode[0] = "if " + conditionCode[0]

	return RawBlock{
		Header:     conditionCode,
		Statements: b.Statements,
	}.CodeLines()
}

func concatenateStatements(sep string, statements []Statement) []string {
	var lines []string
	lastLine := ""
	for _, stmt := range statements {
		stmtLines := stmt.CodeLines()

		if lastLine != "" {
			lastLine += sep
		}
		lastLine += stmtLines[0]

		if len(stmtLines) > 1 {
			lines = append(lines, lastLine)
			lines = append(lines, stmtLines[1:len(stmtLines)-1]...)
			lastLine = stmtLines[len(stmtLines)-1]
		}
	}

	if lastLine != "" {
		lines = append(lines, lastLine)
	}

	return lines
}

type ChainBuilder []Statement

func NewChainBuilder(object string) ChainBuilder {
	return ChainBuilder{
		Identifier(object),
	}
}

func (b ChainBuilder) Chain(field string) ChainBuilder {
	b = append(b, Identifier(field))
	return b
}

func (b ChainBuilder) Call(method string, params ...Statement) ChainBuilder {
	b = append(b, CallStatement{
		FuncName: method,
		Params:   params,
	})
	return b
}

func (b ChainBuilder) Build() ChainStatement {
	return ChainStatement(b)
}
