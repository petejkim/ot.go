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
}
