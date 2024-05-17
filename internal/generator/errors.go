package generator

import "errors"

var (
	ErrStructNotFound    = errors.New("struct not found")
	ErrNotNamedStruct    = errors.New("not a named struct")
	ErrInterfaceNotFound = errors.New("interface not found")
	ErrNotInterface      = errors.New("not an interface")
)
