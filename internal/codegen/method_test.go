package codegen_test

import (
	"bytes"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/testutils"
)

func TestMethodBuilderBuild_IgnoreReceiverNoReturn(t *testing.T) {
	fb := codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{Type: "User"},
		Name:     "Init",
		Params:   nil,
		Returns:  nil,
		Body:     `	db.Init(&User{})`,
	}
	expectedCode := `
func (User) Init() {
	db.Init(&User{})
}
`
	buffer := new(bytes.Buffer)

	err := fb.Impl(buffer)

	if err != nil {
		t.Fatal(err)
	}
	actual := buffer.String()
	if err := testutils.ExpectMultiLineString(
		expectedCode,
		actual,
	); err != nil {
		t.Error(err)
	}
}

func TestMethodBuilderBuild_IgnorePoinerReceiverOneReturn(t *testing.T) {
	fb := codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{
			Type:    "User",
			Pointer: true,
		},
		Name:    "Init",
		Params:  nil,
		Returns: []code.Type{code.TypeError},
		Body:    `	return db.Init(&User{})`,
	}
	expectedCode := `
func (*User) Init() error {
	return db.Init(&User{})
}
`
	buffer := new(bytes.Buffer)

	err := fb.Impl(buffer)

	if err != nil {
		t.Fatal(err)
	}
	actual := buffer.String()
	if err := testutils.ExpectMultiLineString(
		expectedCode,
		actual,
	); err != nil {
		t.Error(err)
	}
}

func TestMethodBuilderBuild_UseReceiverMultiReturn(t *testing.T) {
	fb := codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{
			Name: "u",
			Type: "User",
		},
		Name: "WithAge",
		Params: []code.Param{
			{Name: "age", Type: code.TypeInt},
		},
		Returns: []code.Type{code.SimpleType("User"), code.TypeError},
		Body: `	u.Age = age
	return u`,
	}
	expectedCode := `
func (u User) WithAge(age int) (User, error) {
	u.Age = age
	return u
}
`
	buffer := new(bytes.Buffer)

	err := fb.Impl(buffer)

	if err != nil {
		t.Fatal(err)
	}
	actual := buffer.String()
	if err := testutils.ExpectMultiLineString(
		expectedCode,
		actual,
	); err != nil {
		t.Error(err)
	}
}

func TestMethodBuilderBuild_UsePointerReceiverNoReturn(t *testing.T) {
	fb := codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{
			Name:    "u",
			Type:    "User",
			Pointer: true,
		},
		Name: "SetAge",
		Params: []code.Param{
			{Name: "age", Type: code.TypeInt},
		},
		Returns: nil,
		Body:    `	u.Age = age`,
	}
	expectedCode := `
func (u *User) SetAge(age int) {
	u.Age = age
}
`
	buffer := new(bytes.Buffer)

	err := fb.Impl(buffer)

	if err != nil {
		t.Fatal(err)
	}
	actual := buffer.String()
	if err := testutils.ExpectMultiLineString(
		expectedCode,
		actual,
	); err != nil {
		t.Error(err)
	}
}
