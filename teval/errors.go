package teval

import (
	"errors"
	"fmt"
)

var (
	errEmptyTable = errors.New("empty table")

	errBadExpression  = errors.New("wrong expression format")
	errBadColumnIndex = errors.New("wrong column index format")

	errRecursiveDef   = errors.New("recursive definition")
	errDivisionByZero = errors.New("dividing by zero")
)

type cellError struct {
	col string
	row int
	err error
}

func (e cellError) Error() string {
	return fmt.Sprintf("cell %s%d: %v", e.col, e.row, e.err)
}
