package teval

import (
	"fmt"
	"strconv"
)

type cellReference struct {
	col string
	row int
}

func (c cellReference) Eval(lookup lookupFunc, visited map[cellReference]struct{}) (int, error) {
	if _, ok := visited[c]; ok {
		return 0, errRecursiveDef
	}
	expr, err := lookup(c)
	if err != nil {
		return 0, err
	}
	visited[c] = struct{}{}
	return expr.Eval(lookup, visited)
}

func (c cellReference) String() string {
	return fmt.Sprintf("%s%d", c.col, c.row)
}

type lookupFunc func(c cellReference) (expr, error)

type expr interface {
	Eval(lookup lookupFunc, visited map[cellReference]struct{}) (int, error)
	String() string
}

type integer int

func (i integer) Eval(_ lookupFunc, _ map[cellReference]struct{}) (int, error) {
	return int(i), nil
}

func (i integer) String() string {
	return strconv.Itoa(int(i))
}

type binary struct {
	op          rune
	left, right expr
	evaluated   bool
	value       int
}

func (b *binary) Eval(lookup lookupFunc, visited map[cellReference]struct{}) (int, error) {
	if b.evaluated {
		return b.value, nil
	}

	var val int
	var err error
	l, err := b.left.Eval(lookup, visited)
	if err != nil {
		return 0, err
	}

	r, err := b.right.Eval(lookup, visited)
	if err != nil {
		return 0, err
	}

	switch b.op {
	case '+':
		val = l + r
	case '-':
		val = l - r
	case '*':
		val = l * r
	case '/':
		if r == 0 {
			err = errDivisionByZero
			break
		}
		val = l / r
	}
	if err != nil {
		return 0, err
	}

	b.evaluated = true
	b.value = val
	return val, nil
}

func (b *binary) String() string {
	if b.evaluated {
		return strconv.Itoa(b.value)
	}
	return fmt.Sprintf("=%s%c%s", b.left, b.op, b.right)
}
