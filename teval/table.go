package teval

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// Table represents a matrix of expressions
// with assosiated header and row indexes
type Table struct {
	header      []string
	columnNames []string
	rowIndexes  []int
	cells       [][]expr
	rowLookUp   map[int]int
}

// Read table of specified format
func Read(r io.Reader, sep rune) (*Table, error) {
	input := csv.NewReader(r)
	input.Comma = sep

	// Read header
	header, err := input.Read()
	if err != nil {
		return nil, err
	} else if len(header) < 2 {
		return nil, errEmptyTable
	}

	// Read rows
	rows := [][]string{}
	for {
		row, err := input.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil, errEmptyTable
	}

	// Remove empty cell from header
	columns := header[1:]

	// Convert row indexes to numbers
	rowIndexes := make([]int, 0, len(rows))
	for line, row := range rows {
		ind, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, &csv.ParseError{StartLine: 1, Line: line + 1, Err: errBadColumnIndex}
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
	expressions := make([]expr, 0, len(rows)*len(columns))
	for i, row := range rows {
		for column, cell := range row[1:] {
			expr, err := parseExpression(cell)
			if err != nil {
				parseErrs = append(parseErrs, cellError{columns[column], rowIndexes[i], err})
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

	headerLen := len(columns)
	cells := make([][]expr, len(rows))
	for row := 0; row < len(rows); row++ {
		cells[row] = expressions[row*headerLen : (row+1)*headerLen]
	}

	return &Table{
		header:      header,
		columnNames: columns,
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
func Write(w io.Writer, t *Table, sep rune) error {
	writer := csv.NewWriter(w)
	writer.Comma = sep

	// Write header
	if err := writer.Write(t.header); err != nil {
		return err
	}

	// Write rows
	expr := make([]string, len(t.header))
	for i, row := range t.cells {
		expr[0] = strconv.Itoa(t.rowIndexes[i])
		for j, cell := range row {
			expr[j+1] = cell.String()
		}

		if err := writer.Write(expr); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
