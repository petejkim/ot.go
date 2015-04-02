package selection_test

import (
	"reflect"
	"testing"

	"github.com/nitrous-io/ot.go/ot/operation"
	"github.com/nitrous-io/ot.go/ot/selection"
)

func TestRangeTransform(t *testing.T) {
	r := &selection.Range{5, 9}
	top := operation.New().Retain(10)

	if actual, expected := r, r.Transform(top); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Insert("hello").Retain(10)
	if actual, expected := r.Transform(top), (&selection.Range{10, 14}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(5).Insert("hello").Retain(5)
	if actual, expected := r.Transform(top), (&selection.Range{10, 14}); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(6).Insert("hello").Retain(4)
	if actual, expected := r.Transform(top), (&selection.Range{5, 14}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(9).Insert("hello").Retain(1)
	if actual, expected := r.Transform(top), (&selection.Range{5, 14}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(10).Insert("hello")
	if actual, expected := r.Transform(top), (&selection.Range{5, 9}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Delete(5).Retain(5)
	if actual, expected := r.Transform(top), (&selection.Range{0, 4}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(1).Delete(4).Retain(5)
	if actual, expected := r.Transform(top), (&selection.Range{1, 5}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(1).Delete(5).Retain(4)
	if actual, expected := r.Transform(top), (&selection.Range{1, 4}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(1).Delete(6).Retain(3)
	if actual, expected := r.Transform(top), (&selection.Range{1, 3}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(1).Delete(9)
	if actual, expected := r.Transform(top), (&selection.Range{1, 1}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Delete(10)
	if actual, expected := r.Transform(top), (&selection.Range{0, 0}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = operation.New().Retain(2).Insert("abcd").Delete(3).Retain(1).Delete(2).Insert("abc").Retain(1).Insert("e").Delete(1)
	if actual, expected := r.Transform(top), (&selection.Range{6, 12}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
