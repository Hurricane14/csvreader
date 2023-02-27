package teval

import (
	"reflect"
	"testing"
)

func TestParseArgument(t *testing.T) {
	t.Parallel()

	type test struct {
		description string
		input       string
		expr        expr
		err         error
	}
	tests := []test{
		{
			description: "successful parsing of a cell reference",
			input:       "Cell30",
			expr:        cellReference{"Cell", 30},
			err:         nil,
		},
		{
			description: "successful parsing of a number",
			input:       "30",
			expr:        integer(30),
			err:         nil,
		},
		{
			description: "fail to parse cell reference without row number",
			input:       "Cell",
			expr:        nil,
			err:         errBadColumnIndex,
		},
	}

	for _, test := range tests {
		expr, err := parseArgument(test.input)
		if !reflect.DeepEqual(expr, test.expr) || err != test.err {
			t.Fatalf("%s: wanted: (%s, %s), got: (%s, %s)", test.description, test.expr, test.err, expr, err)
		}
	}
}

func TestParseBinary(t *testing.T) {
	t.Parallel()

	type test struct {
		description string
		input       string
		expr        expr
		err         error
	}
	tests := []test{
		{
			description: "successful parsing of a sum of numbers",
			input:       "=5+8",
			expr:        &binary{op: '+', left: integer(5), right: integer(8)},
			err:         nil,
		},
		{
			description: "successful parsing of a subtraction of numbers",
			input:       "=5-8",
			expr:        &binary{op: '-', left: integer(5), right: integer(8)},
			err:         nil,
		},
		{
			description: "successful parsing of a multiplication of numbers",
			input:       "=5*8",
			expr:        &binary{op: '*', left: integer(5), right: integer(8)},
			err:         nil,
		},
		{
			description: "successful parsing of a division of numbers",
			input:       "=5/8",
			expr:        &binary{op: '/', left: integer(5), right: integer(8)},
			err:         nil,
		},
	}

	for _, test := range tests {
		expr, err := parseBinary(test.input)
		if !reflect.DeepEqual(expr, test.expr) || !reflect.DeepEqual(err, test.err) {
			t.Fatalf("%s: wanted: (%s, %s), got: (%s, %s)", test.description, test.expr, test.err, expr, err)
		}
	}
}
