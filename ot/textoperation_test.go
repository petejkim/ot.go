package ot_test

import (
	"reflect"
	"testing"

	"github.com/petejkim/ot.go/ot"
)

func TestNewTextOperation(t *testing.T) {
	top := ot.NewTextOperation()

	if actual := reflect.TypeOf(top); actual != reflect.TypeOf(&ot.TextOperation{}) {
		t.Fatalf("expected NewTextOperation to return a pointer to ot.TextOperation, got %v", actual)
	}

	if top.Ops == nil {
		t.Errorf("expected Ops not to be nil, got nil")
	}

	if actual := top.BaseLen; actual != 0 {
		t.Errorf("expected BaseLen to be 0, got %d", actual)
	}

	if actual := top.TargetLen; actual != 0 {
		t.Errorf("expected TargetLen to be 0, got %d", actual)
	}
}

func TestRetain(t *testing.T) {
	top := &ot.TextOperation{}

	top.Retain(0)

	if actual, expected := top.BaseLen, 0; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 0; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual := len(top.Ops); actual != 0 {
		t.Errorf("expected empty ops, got ops with length %d", actual)
	}

	top.Retain(2)

	if actual, expected := top.BaseLen, 2; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 2; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*ot.Op{
		&ot.Op{N: 2},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top.Retain(3).Retain(1)

	if actual, expected := top.BaseLen, 6; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 6; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*ot.Op{
		&ot.Op{N: 6},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestDelete(t *testing.T) {
	top := &ot.TextOperation{}

	top.Delete(0)

	if actual, expected := top.BaseLen, 0; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 0; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual := len(top.Ops); actual != 0 {
		t.Errorf("expected empty ops, got ops with length %d", actual)
	}

	top.Delete(2)

	if actual, expected := top.BaseLen, 2; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 0; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*ot.Op{
		&ot.Op{N: -2},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top.Delete(1).Delete(2)

	if actual, expected := top.BaseLen, 5; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 0; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*ot.Op{
		&ot.Op{N: -5},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestInsert(t *testing.T) {
	top := &ot.TextOperation{}

	top.Insert("")

	if actual, expected := top.BaseLen, 0; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 0; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual := len(top.Ops); actual != 0 {
		t.Errorf("expected empty ops, got ops with length %d", actual)
	}

	top.Insert("foo")

	if actual, expected := top.BaseLen, 0; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 3; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*ot.Op{
		&ot.Op{S: "foo"},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top.Insert("bar").Insert("baz")

	if actual, expected := top.BaseLen, 0; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 9; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*ot.Op{
		&ot.Op{S: "foobarbaz"},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestMultipleOps(t *testing.T) {
	top := &ot.TextOperation{}

	top.Retain(1).Delete(2).Delete(1).Retain(2).Retain(3).Delete(1)

	if actual, expected := top.BaseLen, 10; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 6; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*ot.Op{
		&ot.Op{N: 1},
		&ot.Op{N: -3},
		&ot.Op{N: 5},
		&ot.Op{N: -1},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = &ot.TextOperation{}

	top.Retain(1).Insert("foo").Delete(2).Insert("bar").Delete(1).Retain(2).Retain(3).Insert("baz").Delete(1).Delete(3)

	if actual, expected := top.BaseLen, 13; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 15; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*ot.Op{
		&ot.Op{N: 1},
		&ot.Op{S: "foobar"},
		&ot.Op{N: -3},
		&ot.Op{N: 5},
		&ot.Op{S: "baz"},
		&ot.Op{N: -4},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestApply(t *testing.T) {
	top := &ot.TextOperation{}
	top.Retain(3)

	_, err := top.Apply("fo")

	if err != ot.ErrBaseLenMismatch {
		t.Errorf("expected ot.ErrBaseLenMismatch, got %v", err)
	}

	_, err = top.Apply("food")

	if err != ot.ErrBaseLenMismatch {
		t.Errorf("expected ot.ErrBaseLenMismatch, got %v", err)
	}

	s, err := top.Apply("foo")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "foo"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	top = &ot.TextOperation{}
	s, err = top.Insert("bar").Apply("")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "bar"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	top = &ot.TextOperation{}
	s, err = top.Delete(3).Apply("baz")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual := s; actual != "" {
		t.Errorf("expected empty string, got %s", actual)
	}

	top = &ot.TextOperation{}
	s, err = top.Retain(1).Insert("dar").Delete(1).Retain(1).Insert("biz").Apply("fox")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "fdarxbiz"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestTransform(t *testing.T) {
	a := &ot.TextOperation{}
	a.Retain(1)

	b := &ot.TextOperation{}
	b.Retain(2)

	_, _, err := ot.Transform(a, b)

	if err != ot.ErrBaseLenMismatch {
		t.Errorf("expected ot.ErrBaseLenMismatch, got %v", err)
	}

	// apply(apply(S, A), B') = apply(apply(S, B), A')

	s := "She is a girl!!!"
	o := "He was one beautiful man."

	a = &ot.TextOperation{}
	a.Retain(4).Delete(1).Insert("wa").Retain(4).Insert("beautiful ").Retain(4).Delete(3).Insert(".")

	b = &ot.TextOperation{}
	b.Delete(2).Insert("H").Retain(5).Delete(1).Insert("one").Retain(1).Delete(4).Insert("man").Delete(2).Retain(1)

	a1, b1, err := ot.Transform(a, b)

	if err != nil {
		t.Fatalf("expected no error transforming, got %v", err)
	}

	if a1 == nil {
		t.Fatalf("expected non-nil a', got nil")
	}

	if b1 == nil {
		t.Fatalf("expected non-nil b', got nil")
	}

	as, err := a.Apply(s)
	if err != nil {
		t.Fatalf("expected no error applying A, got %v", err)
	}

	at, err := b1.Apply(as)
	if err != nil {
		t.Fatalf("expected no error applying B', got %v", err)
	}

	if actual, expected := at, o; actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}

	bs, err := b.Apply(s)
	if err != nil {
		t.Fatalf("expected no error applying B, got %v", err)
	}

	bt, err := a1.Apply(bs)
	if err != nil {
		t.Fatalf("expected no error applying A', got %v", err)
	}

	if actual, expected := bt, o; actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}
