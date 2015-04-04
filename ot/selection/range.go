package selection

import (
	"unicode/utf8"

	"github.com/nitrous-io/ot.go/ot/operation"
)

type Range struct {
	Anchor int `json:"anchor"`
	Head   int `json:"head"`
}

func (r *Range) Transform(op *operation.Operation) *Range {
	return &Range{transformIndex(r.Anchor, op), transformIndex(r.Head, op)}
}

func transformIndex(i int, op *operation.Operation) int {
	// start cursor at index 0
	j := 0

	for _, op := range op.Ops {
		// if cursor index is greater than i, the rest of the ops are irrelevant
		if j > i {
			break
		}
		if operation.IsRetain(op) {
			// advance cursor
			j += op.N
		} else if operation.IsInsert(op) {
			// insertion increments index. also advance cursor
			i += utf8.RuneCountInString(op.S)
			j += utf8.RuneCountInString(op.S)
		} else if operation.IsDelete(op) {
			// deletion decrements index, but only up to current cursor
			i = max(j, i+op.N) // N is negative
		}
	}

	return i
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
