package codegen_test

import (
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
)

func TestIdentifier(t *testing.T) {
	identifier := codegen.Identifier("user")
	expected := []string{"user"}

	actual := identifier.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestDeclStatement(t *testing.T) {
	stmt := codegen.DeclStatement{
		Name: "arrs",
		Type: code.ArrayType{ContainedType: code.SimpleType("int")},
	}
	expected := []string{"var arrs []int"}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestDeclAssignStatement(t *testing.T) {
	stmt := codegen.DeclAssignStatement{
		Vars: []string{"value", "err"},
		Values: codegen.StatementList{
			codegen.Identifier("1"),
			codegen.Identifier("nil"),
		},
	}
	expected := []string{"value, err := 1, nil"}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestAssignStatement(t *testing.T) {
	stmt := codegen.AssignStatement{
		Vars: []string{"value", "err"},
		Values: codegen.StatementList{
			codegen.Identifier("1"),
			codegen.Identifier("nil"),
		},
	}
	expected := []string{"value, err = 1, nil"}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestReturnStatement(t *testing.T) {
	stmt := codegen.ReturnStatement{
		codegen.Identifier("result"),
		codegen.Identifier("nil"),
	}
	expected := []string{"return result, nil"}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestChainStatement(t *testing.T) {
	stmt := codegen.ChainStatement{
		codegen.Identifier("r"),
		codegen.Identifier("userRepository"),
		codegen.CallStatement{
			FuncName: "Insert",
			Params: codegen.StatementList{
				codegen.StructStatement{
					Type: "User",
					Pairs: []codegen.StructFieldPair{
						{
							Key:   "ID",
							Value: codegen.Identifier("arg0"),
						},
						{
							Key:   "Name",
							Value: codegen.Identifier("arg1"),
						},
					},
				},
			},
		},
		codegen.CallStatement{
			FuncName: "Do",
		},
	}
	expected := []string{
		"r.userRepository.Insert(User{",
		"	ID: arg0,",
		"	Name: arg1,",
		"}).Do()",
	}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestCallStatement(t *testing.T) {
	stmt := codegen.CallStatement{
		FuncName: "FindByID",
		Params: codegen.StatementList{
			codegen.Identifier("ctx"),
			codegen.Identifier("user"),
		},
	}
	expected := []string{"FindByID(ctx, user)"}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestSliceStatement(t *testing.T) {
	stmt := codegen.SliceStatement{
		Type: code.ArrayType{
			ContainedType: code.SimpleType("string"),
		},
		Values: []codegen.Statement{
			codegen.Identifier(`"hello"`),
			codegen.ChainStatement{
				codegen.CallStatement{
					FuncName: "GetUser",
					Params: codegen.StatementList{
						codegen.Identifier("userID"),
					},
				},
				codegen.Identifier("Name"),
			},
		},
	}
	expected := []string{
		"[]string{",
		`	"hello",`,
		`	GetUser(userID).Name,`,
		"}",
	}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestMapStatement(t *testing.T) {
	stmt := codegen.MapStatement{
		Type: "map[string]int",
		Pairs: []codegen.MapPair{
			{
				Key:   "key1",
				Value: codegen.Identifier("value1"),
			},
			{
				Key:   "key2",
				Value: codegen.Identifier("value2"),
			},
		},
	}
	expected := []string{
		"map[string]int{",
		`	"key1": value1,`,
		`	"key2": value2,`,
		"}",
	}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestStructStatement(t *testing.T) {
	stmt := codegen.StructStatement{
		Type: "User",
		Pairs: []codegen.StructFieldPair{
			{
				Key:   "ID",
				Value: codegen.Identifier("arg0"),
			},
			{
				Key:   "Name",
				Value: codegen.Identifier("arg1"),
			},
		},
	}
	expected := []string{
		"User{",
		`	ID: arg0,`,
		`	Name: arg1,`,
		"}",
	}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}

func TestIfBlockStatement(t *testing.T) {
	stmt := codegen.IfBlock{
		Condition: []codegen.Statement{
			codegen.DeclAssignStatement{
				Vars: []string{"err"},
				Values: codegen.StatementList{
					codegen.CallStatement{
						FuncName: "Insert",
						Params: codegen.StatementList{
							codegen.Identifier("ctx"),
							codegen.StructStatement{
								Type: "User",
								Pairs: []codegen.StructFieldPair{
									{
										Key:   "ID",
										Value: codegen.Identifier("id"),
									},
									{
										Key:   "Name",
										Value: codegen.Identifier("name"),
									},
								},
							},
						},
					},
				},
			},
			codegen.RawStatement("err != nil"),
		},
		Statements: []codegen.Statement{
			codegen.ReturnStatement{
				codegen.Identifier("nil"),
				codegen.Identifier("err"),
			},
		},
	}
	expected := []string{
		"if err := Insert(ctx, User{",
		"	ID: id,",
		"	Name: name,",
		"}); err != nil {",
		"	return nil, err",
		"}",
	}

	actual := stmt.CodeLines()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected=%+v actual=%+v", expected, actual)
	}
}
