package teval

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Table represents a matrix of expressions
// with assosiated header and row indexes
type Table struct {
	columnNames []string
	rowIndexes  []int
	cells       [][]expr
	rowLookUp   map[int]int
}

// Read table of specified format
func Read(r io.Reader, sep string) (*Table, error) {
	var line int
	nextLine := func(input *bufio.Scanner) ([]string, error) {
		if !input.Scan() {
			if err := input.Err(); err != nil {
				return nil, err
			}
			return nil, io.EOF
		}
		line++
		return strings.Split(input.Text(), sep), nil
	}
	input := bufio.NewScanner(r)

	// Read header
	header, err := nextLine(input)
	if err != nil {
		return nil, err
	} else if len(header) < 2 {
		return nil, errEmptyTable
	}

	// Read rows
	rows := [][]string{}
	for {
		row, err := nextLine(input)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, lineError{line, err}
		}
		if len(row) != len(header) {
			return nil, lineError{line, fmt.Errorf("mismatched number of columns")}
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil, lineError{line, errEmptyTable}
	}

	// Remove empty cell from header
	header = header[1:]

	// Convert row indexes to numbers
	rowIndexes := make([]int, 0, len(rows))
	for _, row := range rows {
		ind, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, lineError{line, fmt.Errorf("row index is not a number")}
		}
		rowIndexes = append(rowIndexes, ind)
	}

	// Create lookup table for row indexes
	rowLookUp := make(map[int]int, len(rowIndexes))
	for ind, row := range rowIndexes {
		rowLookUp[row] = ind
	}

	// Parse expressions
	parseErrs := []error{}
	expressions := make([]expr, 0, len(rows)*len(header))
	for i, row := range rows {
		for column, cell := range row[1:] {
			expr, err := parseExpression(cell)
			if err != nil {
				parseErrs = append(parseErrs, cellError{header[column], rowIndexes[i], err})
			}
			if len(parseErrs) > 10 {
				parseErrs = append(parseErrs, fmt.Errorf("too many parsing errors, parsing stopped..."))
				break
			}
			expressions = append(expressions, expr)
		}
	}
	if len(parseErrs) != 0 {
		return nil, errors.Join(parseErrs...)
	}

	headerLen := len(header)
	cells := make([][]expr, len(rows))
	for row := 0; row < len(rows); row++ {
		cells[row] = expressions[row*headerLen : (row+1)*headerLen]
	}

	return &Table{
		columnNames: header,
		rowIndexes:  rowIndexes,
		cells:       cells,
		rowLookUp:   rowLookUp,
	}, nil
}

// Returns expression in specified cell
func (t *Table) lookup(cr cellReference) (expr, error) {
	var col, row int

	col = -1
	for i, column := range t.columnNames {
		if column == cr.col {
			col = i
			break
		}
	}
	if col < 0 {
		return nil, fmt.Errorf("column %s does not exist", cr.col)
	}

	row, ok := t.rowLookUp[cr.row]
	if !ok {
		return nil, fmt.Errorf("row %d does not exist", cr.row)
	}

	return t.cells[row][col], nil
}

// Evaluates all cells of the table
func (t *Table) EvalAll() error {
	errs := []error{}
	for row := range t.cells {
		for column := range t.cells[row] {
			cr := cellReference{t.columnNames[column], t.rowIndexes[row]}
			cell := t.cells[row][column]
			_, err := cell.Eval(t.lookup, map[cellReference]struct{}{cr: {}})
			if err != nil {
				errs = append(errs, cellError{t.columnNames[column], t.rowIndexes[row], err})
			}
			if len(errs) > 10 {
				errs = append(errs, fmt.Errorf("too many errors, evaluation stopped..."))
				break
			}
		}
	}
	return errors.Join(errs...)
}

// Writes table in the specified format using sep as separator
// If table was evaluated, results of evaluation are used instead of expressions
func Write(w io.Writer, t *Table, sep string) error {
	var err error
	bw := bufio.NewWriter(w)

	// Write header
	_, err = fmt.Fprintf(bw, "%s%s\n", sep, strings.Join(t.columnNames, sep))
	if err != nil {
		return err
	}

	// Write rows
	expr := make([]string, len(t.columnNames))
	for i, row := range t.cells {
		for j, cell := range row {
			expr[j] = cell.String()
		}

		format := "%d%s%s\n"
		if i == len(t.cells)-1 {
			format = format[:len(format)-1]
		}

		_, err := fmt.Fprintf(bw, format, t.rowIndexes[i], sep, strings.Join(expr, sep))
		if err != nil {
			return err
		}
	}

	return bw.Flush()
}
