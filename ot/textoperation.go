package ot

import (
	"errors"
	"fmt"
)

var (
	ErrBaseLenMismatch = errors.New("ot: base length mismatch")
	ErrTransformFailed = errors.New("ot: transform failed")
	ErrMarshalFailed   = errors.New("ot: marshal failed")
	ErrUnmarshalFailed = errors.New("ot: unmarshal failed")
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

func NewTextOperation() *TextOperation {
	return &TextOperation{Ops: []*Op{}}
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

func (t *TextOperation) At(i int) *Op {
	if i >= len(t.Ops) {
		return nil
	}
	return t.Ops[i]
}

func (t *TextOperation) Marshal() []interface{} {
	ops := make([]interface{}, len(t.Ops))

	for i, o := range t.Ops {
		if o.S == "" {
			ops[i] = o.N
		} else {
			ops[i] = o.S
		}
	}

	return ops
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

func Transform(a, b *TextOperation) (*TextOperation, *TextOperation, error) {
	if a.BaseLen != b.BaseLen {
		return nil, nil, ErrBaseLenMismatch
	}

	a1, b1 := &TextOperation{}, &TextOperation{}
	iA, iB := 0, 0
	opA, opB := a.At(iA), b.At(iB)

	nextOpA := func() {
		iA++
		opA = a.At(iA)
	}
	nextOpB := func() {
		iB++
		opB = b.At(iB)
	}

	for !(opA == nil && opB == nil) {
		// either op is insert e.g. Op A=insert => A'<- insert, B'<- retain
		// if both are insert, process op A first
		if opA != nil && IsInsert(opA) {
			a1.Insert(opA.S)
			b1.Retain(len(opA.S))
			nextOpA()
			continue
		} else if opB != nil && IsInsert(opB) {
			a1.Retain(len(opB.S))
			b1.Insert(opB.S)
			nextOpB()
			continue
		}

		if opA == nil || opB == nil {
			return nil, nil, ErrTransformFailed
		}

		// retain/retain
		if IsRetain(opA) && IsRetain(opB) {
			min, nA, nB := 0, opA.N, opB.N
			if nA > nB {
				min = nB
				opA = &Op{N: nA - nB}
				nextOpB()
			} else if nA < nB {
				min = nA
				nextOpA()
				opB = &Op{N: nB - nA}
			} else {
				min = nA
				nextOpA()
				nextOpB()
			}
			a1.Retain(min)
			b1.Retain(min)
			continue
		}

		// delete/delete
		// both deleting at same index, handle where one deletes more
		if IsDelete(opA) && IsDelete(opB) {
			nA, nB := -opA.N, -opB.N
			if nA > nB {
				opA = &Op{N: -(nA - nB)}
				nextOpB()
			} else if nA < nB {
				nextOpA()
				opB = &Op{N: -(nB - nA)}
			} else {
				nextOpA()
				nextOpB()
			}
			continue
		}

		// delete/retain
		if IsDelete(opA) && IsRetain(opB) {
			min, nA, nB := 0, -opA.N, opB.N
			if nA > nB {
				min = nB
				opA = &Op{N: nB - nA} // delete
				nextOpB()
			} else if nA < nB {
				min = nA
				nextOpA()
				opB = &Op{N: nB - nA} // retain
			} else {
				min = nA
				nextOpA()
				nextOpB()
			}
			a1.Delete(min)
			continue
		}

		// retain/delete
		if IsRetain(opA) && IsDelete(opB) {
			min, nA, nB := 0, opA.N, -opB.N
			if nA > nB {
				min = nB
				opA = &Op{N: nA - nB} // retain
				nextOpB()
			} else if nA < nB {
				min = nA
				nextOpA()
				opB = &Op{N: nA - nB} // delete
			} else {
				min = nA
				nextOpA()
				nextOpB()
			}
			b1.Delete(min)
			continue
		}

		return nil, nil, ErrTransformFailed
	}

	return a1, b1, nil
}

func Unmarshal(ops []interface{}) (*TextOperation, error) {
	top := &TextOperation{}
	for _, o := range ops {
		switch o.(type) {
		case int:
			n := o.(int)
			if n > 0 {
				top.Retain(n)
			} else {
				top.Delete(-n)
			}
		case float64:
			n := int(o.(float64))
			if n > 0 {
				top.Retain(n)
			} else {
				top.Delete(-n)
			}
		case string:
			s := o.(string)
			top.Insert(s)
		default:
			return nil, ErrUnmarshalFailed
		}
	}
	return top, nil
}
