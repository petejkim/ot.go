package selection

import (
	"errors"

	"github.com/nitrous-io/ot.go/ot/operation"
)

var (
	ErrUnmarshalFailed = errors.New("ot/selection: unmarshal failed")
)

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

func (s *Selection) Marshal() map[string]interface{} {
	mr := make([]map[string]interface{}, len(s.Ranges))
	for i, r := range s.Ranges {
		mr[i] = map[string]interface{}{"anchor": r.Anchor, "head": r.Head}
	}
	return map[string]interface{}{"ranges": mr}
}

func Unmarshal(data map[string]interface{}) (*Selection, error) {
	if data["ranges"] == nil {
		return nil, ErrUnmarshalFailed
	}

	dr, ok := data["ranges"].([]interface{})
	if !ok {
		return nil, ErrUnmarshalFailed
	}

	ranges := make([]Range, len(dr))

	for i, o := range dr {
		r, ok := o.(map[string]interface{})
		if !ok {
			return nil, ErrUnmarshalFailed
		}
		rng, err := unmarshalRange(r)
		if err != nil {
			return nil, err
		}
		ranges[i] = *rng
	}

	return &Selection{ranges}, nil
}

func unmarshalRange(data map[string]interface{}) (*Range, error) {
	a, ok := parseNumber(data["anchor"])
	if !ok {
		return nil, ErrUnmarshalFailed
	}
	h, ok := parseNumber(data["head"])
	if !ok {
		return nil, ErrUnmarshalFailed
	}
	return &Range{a, h}, nil
}

func parseNumber(n interface{}) (int, bool) {
	switch n.(type) {
	case int:
		return n.(int), true
	case float64:
		return int(n.(float64)), true
	}
	return 0, false
}
