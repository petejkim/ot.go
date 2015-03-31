package ot

import (
	"errors"
	"fmt"
)

var (
	ErrBaseLenMismatch = errors.New("ot: base length mismatch")
)

type Op struct {
	N int
	S string
}

func (o *Op) String() string {
	return fmt.Sprintf("&%+v", *o)
}

type TextOperation struct {
	Ops       []*Op
	BaseLen   int
	TargetLen int
}

func (t *TextOperation) Retain(n int) *TextOperation {
	if n <= 0 {
		return t
	}
	t.BaseLen += n
	t.TargetLen += n

	last := t.LastOp()
	if last != nil && IsRetain(last) {
		// last op is retain -> merge
		last.N += n
	} else {
		t.Ops = append(t.Ops, &Op{N: n})
	}
	return t
}

func (t *TextOperation) Delete(n int) *TextOperation {
	if n <= 0 {
		return t
	}
	t.BaseLen += n

	last := t.LastOp()
	if last != nil && IsDelete(last) {
		// last op is delete -> merge
		last.N -= n
	} else {
		t.Ops = append(t.Ops, &Op{N: -n})
	}

	return t
}

func (t *TextOperation) Insert(s string) *TextOperation {
	if s == "" {
		return t
	}
	t.TargetLen += len(s)

	last := t.LastOp()
	if last != nil && IsInsert(last) {
		// last op is insert -> merge
		last.S += s
	} else if last != nil && IsDelete(last) {
		// last op is delete -> put insert before the delete
		var secondLast *Op
		opsLen := len(t.Ops)
		if opsLen >= 2 {
			secondLast = t.Ops[opsLen-2]
		}
		if secondLast != nil && IsInsert(secondLast) {
			// 2nd last op is insert -> merge
			secondLast.S += s
		} else {
			t.Ops = append(t.Ops, last)
			t.Ops[opsLen-1] = &Op{S: s}
		}
	} else {
		t.Ops = append(t.Ops, &Op{S: s})
	}

	return t
}

func (t *TextOperation) LastOp() *Op {
	if len(t.Ops) == 0 {
		return nil
	}
	return t.Ops[len(t.Ops)-1]
}

func (t *TextOperation) Apply(s string) (string, error) {
	if len(s) != t.BaseLen {
		return "", ErrBaseLenMismatch
	}

	newStr := ""
	// start cursor at index 0 of original string
	i := 0

	for _, op := range t.Ops {
		if IsRetain(op) {
			// copy retained chars and advance cursor
			newStr += s[i : i+op.N]
			i += op.N
		} else if IsInsert(op) {
			// copy inserted chars, but do not advance cursor
			newStr += op.S
		} else if IsDelete(op) {
			// skip deleted chars by advancing cursor
			i -= op.N // N is negative
		}
	}

	return newStr, nil
}

func IsRetain(op *Op) bool {
	return op.N > 0 && op.S == ""
}

func IsDelete(op *Op) bool {
	return op.N < 0 && op.S == ""
}

func IsInsert(op *Op) bool {
	return op.N == 0 && op.S != ""
}
