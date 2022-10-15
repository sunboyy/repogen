package code

import (
	"errors"
	"fmt"
)

var (
	ErrAmbiguousPackageName = errors.New("code: ambiguous package name")
)

type DuplicateStructError string

func (err DuplicateStructError) Error() string {
	return fmt.Sprintf(
		"code: duplicate implementation of struct '%s'",
		string(err),
	)
}

type DuplicateInterfaceError string

func (err DuplicateInterfaceError) Error() string {
	return fmt.Sprintf(
		"code: duplicate implementation of interface '%s'",
		string(err),
	)
}
