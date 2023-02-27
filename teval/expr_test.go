package teval

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEvalBinary(t *testing.T) {
	t.Parallel()

	type test struct {
		description string
		expr        binary
		lookup      lookupFunc
		visited     map[cellReference]struct{}
		value       int
		err         error
	}
	tests := []test{
		{
			description: "sum of two numbers",
			expr:        binary{op: '+', left: integer(3), right: integer(3)},
			value:       6,
			err:         nil,
		},
		{
			description: "difference of two numbers",
			expr:        binary{op: '-', left: integer(3), right: integer(3)},
			value:       0,
			err:         nil,
		},
		{
			description: "multiplication of two numbers",
			expr:        binary{op: '*', left: integer(3), right: integer(3)},
			value:       9,
			err:         nil,
		},
		{
			description: "division of two numbers",
			expr:        binary{op: '/', left: integer(3), right: integer(3)},
			value:       1,
			err:         nil,
		},
		{
			description: "division by zero",
			expr:        binary{op: '/', left: integer(3), right: integer(0)},
			value:       0,
			err:         errDivisionByZero,
		},
		{
			description: "sum of cell and number",
			expr:        binary{op: '+', left: cellReference{}, right: integer(1)},
			lookup:      func(c cellReference) (expr, error) { return integer(3), nil },
			visited:     map[cellReference]struct{}{},
			value:       4,
			err:         nil,
		},
	}

	for _, test := range tests {
		val, err := test.expr.Eval(test.lookup, test.visited)
		if test.value != val || !reflect.DeepEqual(err, test.err) {
			t.Fatalf("%s: wanted: (%d, %v), got: (%d, %v)", test.description, test.value, test.err, val, err)
		}
	}
}

func TestEvalCellReference(t *testing.T) {
	t.Parallel()

	type test struct {
		description string
		ref         cellReference
		lookup      lookupFunc
		visited     map[cellReference]struct{}
		value       int
		err         error
	}
	tests := []test{
		{
			description: "lookup receives unmodified values",
			ref:         cellReference{"col", 1},
			lookup: func(c cellReference) (expr, error) {
				if c.col != "col" || c.row != 1 {
					return nil, fmt.Errorf("unexpected cell reference")
				}
				return integer(3), nil
			},
			visited: map[cellReference]struct{}{},
			value:   3,
			err:     nil,
		},
		{
			description: "recursive references are not allowed",
			ref:         cellReference{"col", 1},
			lookup:      nil,
			visited:     map[cellReference]struct{}{{"col", 1}: {}},
			value:       0,
			err:         errRecursiveDef,
		},
	}

	for _, test := range tests {
		val, err := test.ref.Eval(test.lookup, test.visited)
		if test.value != val || !reflect.DeepEqual(err, test.err) {
			t.Fatalf("%s: wanted: (%d, %v), got: (%d, %v)", test.description, test.value, test.err, val, err)
		}
	}
}
