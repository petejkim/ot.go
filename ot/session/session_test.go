package session_test

import (
	"reflect"
	"testing"

	"github.com/nitrous-io/ot.go/ot/operation"
	"github.com/nitrous-io/ot.go/ot/session"
)

func TestNew(t *testing.T) {
	doc := "Lorem Ipsum Dolor Sit Amet"
	s := session.New(doc)

	if actual := reflect.TypeOf(s); actual != reflect.TypeOf(&session.Session{}) {
		t.Fatalf("expected NewSession to return a pointer to session.Session, got %v", actual)
	}

	if actual := s.Document; actual != doc {
		t.Errorf("expected document to be %s, got %s", doc, actual)
	}

	if s.Operations == nil {
		t.Errorf("expected operations not to be nil, got nil")
	}

	if s.Clients == nil {
		t.Errorf("expected clients not to be nil, got nil")
	}
}

func TestAddClient(t *testing.T) {
	s := session.New("")
	s.AddClient("foo")

	cl := s.Clients["foo"]

	if cl == nil {
		t.Fatalf("expected client to have been added, but it was not")
	}

	if actual := cl.Name; actual != "" {
		t.Errorf("expected name to be empty, but got %s", actual)
	}

	if actual := cl.Selection.Ranges; actual == nil {
		t.Errorf("expected selection ranges not to be nil, got nil")
	}

	if actual := cl.Selection.Ranges; len(actual) != 0 {
		t.Errorf("expected selection ranges to be empty, but got %+v", actual)
	}
}

func TestRemoveClient(t *testing.T) {
	s := session.New("")
	s.AddClient("foo")
	s.AddClient("bar")
	s.AddClient("baz")
	s.RemoveClient("bar")

	if s.Clients["foo"] == nil {
		t.Errorf("expected foo not to have been removed, but it was")
	}

	if s.Clients["bar"] != nil {
		t.Errorf("expected bar to have been removed, but it was not")
	}

	if s.Clients["baz"] == nil {
		t.Errorf("expected baz not to have been removed, but it was")
	}
}

func TestAddOperation(t *testing.T) {
	s := session.New("I love you.")

	// I love you. -> She love you.
	op1 := operation.New().Delete(1).Insert("She").Retain(10)
	// She love you. -> She loves you.
	op2 := operation.New().Retain(8).Insert("s").Retain(5)
	// She loves you. -> She loves you!!!
	op3 := operation.New().Retain(13).Insert("!!!").Delete(1)

	// invalid revision (given rev(1) > current rev(0))
	_, err := s.AddOperation(1, op1)

	if err != session.ErrInvalidRevision {
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
	op4 := operation.New().Retain(2).Insert("really ").Retain(9)

	retOp, err = s.AddOperation(0, op4)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// returned op should be transformed against concurrent operations
	if expected := operation.New().Retain(4).Insert("really ").Retain(12); !reflect.DeepEqual(retOp, expected) {
		t.Errorf("expected returned operation to equal %v, got %v", expected, retOp)
	}

	// rollback last operation
	s.Operations = s.Operations[:len(s.Operations)-1]
	s.Document = prevDoc

	// simulate client that is at revision 1
	op5 := operation.New().Retain(4).Insert("really ").Retain(9)
	retOp, err = s.AddOperation(1, op5)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// returned op should be transformed against concurrent operations
	if expected := operation.New().Retain(4).Insert("really ").Retain(12); !reflect.DeepEqual(retOp, expected) {
		t.Errorf("expected returned operation to equal %v, got %v", expected, retOp)
	}
}
