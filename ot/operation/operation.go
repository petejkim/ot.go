package operation

import (
	"errors"
	"fmt"
	"unicode/utf16"

	"github.com/nitrous-io/ot.go/ot"
)

var (
	ErrBaseLenMismatch = errors.New("ot/operation: base length mismatch")
	ErrTransformFailed = errors.New("ot/operation: transform failed")
	ErrMarshalFailed   = errors.New("ot/operation: marshal failed")
	ErrUnmarshalFailed = errors.New("ot/operation: unmarshal failed")
)

type Op struct {
	N int
	S []rune
}

func (o *Op) String() string {
	return fmt.Sprintf("&%+v", *o)
}

type Operation struct {
	Ops       []*Op
	BaseLen   int
	TargetLen int
	Meta      interface{}
}

func New() *Operation {
	return &Operation{Ops: []*Op{}}
}

func (t *Operation) Retain(n int) *Operation {
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

func (t *Operation) Delete(n int) *Operation {
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

func (t *Operation) Insert(s string) *Operation {
	if s == "" {
		return t
	}

	r := []rune(s)
	if ot.TextEncoding == ot.TextEncodingTypeUTF16 {
		r = uint16sToRunes(utf16.Encode(r))
	}
	t.TargetLen += len(r)

	last := t.LastOp()
	if last != nil && IsInsert(last) {
		// last op is insert -> merge
		last.S = append(last.S, r...)
	} else if last != nil && IsDelete(last) {
		// last op is delete -> put insert before the delete
		var secondLast *Op
		opsLen := len(t.Ops)
		if opsLen >= 2 {
			secondLast = t.Ops[opsLen-2]
		}
		if secondLast != nil && IsInsert(secondLast) {
			// 2nd last op is insert -> merge
			secondLast.S = append(secondLast.S, r...)
		} else {
			t.Ops = append(t.Ops, last)
			t.Ops[opsLen-1] = &Op{S: r}
		}
	} else {
		t.Ops = append(t.Ops, &Op{S: r})
	}

	return t
}

func (t *Operation) LastOp() *Op {
	if len(t.Ops) == 0 {
		return nil
	}
	return t.Ops[len(t.Ops)-1]
}

func (t *Operation) Apply(s string) (string, error) {
	r := []rune(s)
	if ot.TextEncoding == ot.TextEncodingTypeUTF16 {
		r = uint16sToRunes(utf16.Encode(r))
	}

	if len(r) != t.BaseLen {
		return "", ErrBaseLenMismatch
	}

	newStr := ""
	// start cursor at index 0 of original string
	i := 0

	for _, op := range t.Ops {
		if IsRetain(op) {
			// copy retained chars and advance cursor
			ss := r[i : i+op.N]
			if ot.TextEncoding == ot.TextEncodingTypeUTF8 {
				newStr += string(ss)
			} else {
				newStr += string(utf16.Decode(runesToUint16s(ss)))
			}
			i += op.N
		} else if IsInsert(op) {
			// copy inserted chars, but do not advance cursor
			if ot.TextEncoding == ot.TextEncodingTypeUTF8 {
				newStr += string(op.S)
			} else {
				newStr += string(utf16.Decode(runesToUint16s(op.S)))
			}
		} else if IsDelete(op) {
			// skip deleted chars by advancing cursor
			i -= op.N // N is negative
		}
	}

	return newStr, nil
}

func (t *Operation) At(i int) *Op {
	if i >= len(t.Ops) {
		return nil
	}
	return t.Ops[i]
}

func (t *Operation) Marshal() []interface{} {
	ops := make([]interface{}, len(t.Ops))

	for i, o := range t.Ops {
		if IsInsert(o) {
			if ot.TextEncoding == ot.TextEncodingTypeUTF8 {
				ops[i] = string(o.S)
			} else {
				ops[i] = string(utf16.Decode(runesToUint16s(o.S)))
			}
		} else {
			ops[i] = o.N
		}
	}

	return ops
}

func IsRetain(op *Op) bool {
	return op.N > 0
}

func IsDelete(op *Op) bool {
	return op.N < 0
}

func IsInsert(op *Op) bool {
	return op.N == 0 && op.S != nil && len(op.S) != 0
}

func Transform(a, b *Operation) (*Operation, *Operation, error) {
	if a.BaseLen != b.BaseLen {
		return nil, nil, ErrBaseLenMismatch
	}

	a1, b1 := &Operation{}, &Operation{}
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
			if ot.TextEncoding == ot.TextEncodingTypeUTF8 {
				a1.Insert(string(opA.S))
			} else {
				a1.Insert(string(utf16.Decode(runesToUint16s(opA.S))))
			}
			b1.Retain(len(opA.S))
			nextOpA()
			continue
		} else if opB != nil && IsInsert(opB) {
			a1.Retain(len(opB.S))
			if ot.TextEncoding == ot.TextEncodingTypeUTF8 {
				b1.Insert(string(opB.S))
			} else {
				b1.Insert(string(utf16.Decode(runesToUint16s(opB.S))))
			}
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

func Unmarshal(ops []interface{}) (*Operation, error) {
	top := &Operation{}
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

func uint16sToRunes(s []uint16) []rune {
	r := make([]rune, len(s))
	for i, v := range s {
		r[i] = rune(v)
	}
	return r
}

func runesToUint16s(r []rune) []uint16 {
	s := make([]uint16, len(r))
	for i, v := range r {
		s[i] = uint16(v)
	}
	return s
}
