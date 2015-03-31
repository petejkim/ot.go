package ot_test

import (
	"reflect"
	"testing"

	"github.com/petejkim/ot.go/ot"
)

func TestNewSession(t *testing.T) {
	doc := "Lorem Ipsum Dolor Sit Amet"
	s := ot.NewSession(doc)

	if actual := reflect.TypeOf(s); actual != reflect.TypeOf(&ot.Session{}) {
		t.Fatalf("expected NewSession to return a pointer to ot.Session, got %v", actual)
	}

	if actual := s.Document; actual != doc {
		t.Errorf("expected document to be %s, got %s", doc, actual)
	}

	if s.Operations == nil {
		t.Errorf("expected operations not to be nil, got nil")
	}
}

func TestAddOperation(t *testing.T) {
	s := ot.NewSession("I love you.")

	// I love you. -> She love you.
	op1 := ot.NewTextOperation().Delete(1).Insert("She").Retain(10)
	// She love you. -> She loves you.
	op2 := ot.NewTextOperation().Retain(8).Insert("s").Retain(5)
	// She loves you. -> She loves you!!!
	op3 := ot.NewTextOperation().Retain(13).Insert("!!!").Delete(1)

	// invalid revision (given rev(1) > current rev(0))
	_, err := s.AddOperation(1, op1)

	if err != ot.ErrInvalidRevision {
		t.Errorf("expected ErrInvalidRevision, got %v", err)
	}

	// adding operations with no concurrent operations to transform against
	retOp, err := s.AddOperation(0, op1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(retOp, op1) {
		t.Errorf("expected returned operation to equal %v, got %v", op1, retOp)
	}

	retOp, err = s.AddOperation(1, op2)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(retOp, op2) {
		t.Errorf("expected returned operation to equal %v, got %v", op2, retOp)
	}

	retOp, err = s.AddOperation(2, op3)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(retOp, op3) {
		t.Errorf("expected returned operation to equal %v, got %v", op3, retOp)
	}

	prevDoc := s.Document

	// simulate client that is still at revision 0
	op4 := ot.NewTextOperation().Retain(2).Insert("really ").Retain(9)

	retOp, err = s.AddOperation(0, op4)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// returned op should be transformed against concurrent operations
	if expected := ot.NewTextOperation().Retain(4).Insert("really ").Retain(12); !reflect.DeepEqual(retOp, expected) {
		t.Errorf("expected returned operation to equal %v, got %v", expected, retOp)
	}

	// rollback last operation
	s.Operations = s.Operations[:len(s.Operations)-1]
	s.Document = prevDoc

	// simulate client that is at revision 1
	op5 := ot.NewTextOperation().Retain(4).Insert("really ").Retain(9)
	retOp, err = s.AddOperation(1, op5)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// returned op should be transformed against concurrent operations
	if expected := ot.NewTextOperation().Retain(4).Insert("really ").Retain(12); !reflect.DeepEqual(retOp, expected) {
		t.Errorf("expected returned operation to equal %v, got %v", expected, retOp)
	}
}
