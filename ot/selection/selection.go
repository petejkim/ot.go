package selection

import "github.com/nitrous-io/ot.go/ot/operation"

type Selection struct {
	Ranges []Range `json:"ranges"`
}

func (s *Selection) Transform(op *operation.Operation) *Selection {
	tr := make([]Range, len(s.Ranges))
	for i, r := range s.Ranges {
		tr[i] = *r.Transform(op)
	}
	return &Selection{tr}
}
