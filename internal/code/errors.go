package code

import "fmt"

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
