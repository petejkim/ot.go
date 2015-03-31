package ot

import "errors"

var (
	ErrInvalidRevision = errors.New("ot: invalid revision")
)

type Session struct {
	Document   string
	Operations []*TextOperation
}

func NewSession(document string) *Session {
	return &Session{
		Document:   document,
		Operations: []*TextOperation{},
	}
}

func (s *Session) AddOperation(revision int, operation *TextOperation) (*TextOperation, error) {
	if revision < 0 || len(s.Operations) < revision {
		return nil, ErrInvalidRevision
	}
	// find concurrent operations client isn't yet aware of
	otherOps := s.Operations[revision:]

	// transform given operation against these operations
	for _, otherOp := range otherOps {
		op1, _, err := Transform(operation, otherOp)
		if err != nil {
			return nil, err
		}
		operation = op1
	}

	// apply transformed op on the doc
	doc, err := operation.Apply(s.Document)
	if err != nil {
		return nil, err
	}

	s.Document = doc
	s.Operations = append(s.Operations, operation)

	return operation, nil
}
