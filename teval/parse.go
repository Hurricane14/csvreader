package teval

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Parse a number or cell reference
func parseArgument(s string) (expr, error) {
	if v, err := strconv.Atoi(s); err == nil {
		return integer(v), nil
	}

	var columnName int
	for {
		r, s := utf8.DecodeRuneInString(s[columnName:])
		if r == 0 || !unicode.IsLetter(r) {
			break
		} else if r == utf8.RuneError {
			return nil, fmt.Errorf("couldn't decode utf8 string")
		}
		columnName += s
	}
	if columnName == 0 || columnName == len(s) {
		return nil, errBadColumnIndex
	}
	column := s[:columnName]
	row, err := strconv.Atoi(s[columnName:])
	if err != nil {
		return nil, errBadColumnIndex
	}
	return cellReference{column, row}, nil
}

// Parses "= arg1 op arg2"
func parseBinary(s string) (expr, error) {
	const ops = "+-*/"

	if !strings.HasPrefix(s, "=") {
		return nil, errBadExpression
	}
	s = s[1:]

	opInd := strings.IndexAny(s, ops)
	if opInd < 0 || opInd == len(s) {
		return nil, errBadExpression
	}
	/*
		for i := 0; i < len(ops); i++ {
			op = ops[i]
			args = strings.Split(s, string(op))
			if len(args) == 2 {
				break
			}
		}
		if len(args) != 2 {
			return nil, errBadExpression
		}
	*/

	arg1, err := parseArgument(s[:opInd])
	if err != nil {
		return nil, err
	}
	arg2, err := parseArgument(s[opInd+1:])
	if err != nil {
		return nil, err
	}
	return &binary{
		op:   rune(s[opInd]),
		left: arg1, right: arg2,
	}, nil
}

func parseExpression(s string) (expr, error) {
	if v, err := strconv.Atoi(s); err == nil {
		return integer(v), nil
	}
	return parseBinary(s)
}
