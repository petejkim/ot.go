package ot_test

import (
	"reflect"
	"testing"

	"github.com/petejkim/ot.go/ot"
)

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
