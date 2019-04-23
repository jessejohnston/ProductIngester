package parser

import (
	"fmt"
)

// Error is a product parser error
type Error struct {
	line  int
	col   int
	field []byte
	msg   string
	err   error
}

// NewError creates a new product parser error
func NewError(line, col int, field []byte, msg string, err error) Error {
	return Error{
		line:  line,
		col:   col,
		field: field,
		msg:   msg,
		err:   err,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d, %d): \"%s\" %s: %v", e.line, e.col, string(e.field), e.msg, e.err)
}
