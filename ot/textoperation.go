package ot

import "fmt"

type Op struct {
	N int
	S string
}

var NoOp = Op{}

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
		// insert op
		t.Ops = append(t.Ops, &Op{N: n})
	}
	return t
}

func (t *TextOperation) LastOp() *Op {
	if len(t.Ops) == 0 {
		return nil
	}
	return t.Ops[len(t.Ops)-1]
}

func IsRetain(op *Op) bool {
	return op.N > 0 && op.S == ""
}
