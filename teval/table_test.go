package teval

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	t.Parallel()

	type test struct {
		description string
		data        string
		sep         string
		err         error
	}
	tests := []test{
		{
			description: "successful reading of table",
			data:        ",A,B,Cell\n1,2,0,1\n2,=A1/1,=A1+Cell30,0\n30,0,=B1+A1,5",
			sep:         ",",
			err:         nil,
		},
		{
			description: "empty table",
			data:        ",A,B,Cell\n",
			sep:         ",",
			err:         lineError{1, errEmptyTable},
		},
		{
			description: "empty table",
			data:        ",A,B,Cell",
			sep:         ",",
			err:         lineError{1, errEmptyTable},
		},
		{
			description: "mismatched number of columns",
			data:        ",A,B,Cell\n1,2,0\n",
			sep:         ",",
			err:         lineError{2, fmt.Errorf("mismatched number of columns")},
		},
		{
			description: "row index is not a number",
			data:        ",A,B,Cell\nrow,2,0,0\n",
			sep:         ",",
			err:         lineError{2, fmt.Errorf("row index is not a number")},
		},
	}

	for _, test := range tests {
		r := bytes.NewReader([]byte(test.data))
		_, err := Read(r, test.sep)
		if !reflect.DeepEqual(test.err, err) {
			t.Fatalf("%s: want: %v, got: %v", test.description, test.err, err)
		}
	}
}
