package teval

import (
	"bytes"
	"encoding/csv"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	t.Parallel()

	const sep = ','
	type test struct {
		description string
		data        string
		err         error
	}
	tests := []test{
		{
			description: "successful reading of table",
			data:        ",A,B,Cell\n1,2,0,1\n2,=A1/1,=A1+Cell30,0\n30,0,=B1+A1,5",
			err:         nil,
		},
		{
			description: "empty table with new line",
			data:        ",A,B,Cell\n",
			err:         errEmptyTable,
		},
		{
			description: "empty table",
			data:        ",A,B,Cell",
			err:         errEmptyTable,
		},
		{
			description: "mismatched number of columns",
			data:        ",A,B,Cell\n1,2,0\n",
			err:         &csv.ParseError{StartLine: 2, Line: 2, Column: 1, Err: csv.ErrFieldCount},
		},
		{
			description: "row index is not a number",
			data:        ",A,B,Cell\nrow,2,0,0\n",
			err:         &csv.ParseError{StartLine: 1, Line: 1, Column: 0, Err: errBadColumnIndex},
		},
	}

	for _, test := range tests {
		r := bytes.NewReader([]byte(test.data))
		_, err := Read(r, sep)
		if !reflect.DeepEqual(test.err, err) {
			t.Fatalf("%s: want: %v, got: %v", test.description, test.err, err)
		}
	}
}
