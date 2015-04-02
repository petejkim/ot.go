package selection_test

import (
	"reflect"
	"testing"

	"github.com/nitrous-io/ot.go/ot/operation"
	"github.com/nitrous-io/ot.go/ot/selection"
)

func TestSelectionTransform(t *testing.T) {
	s := &selection.Selection{[]selection.Range{{5, 8}, {2, 11}}}

	top := operation.New().Retain(1).Delete(2).Insert("Hello").Retain(4).Delete(2).Insert("Woo").Retain(4)

	if actual, expected := s.Transform(top), (&selection.Selection{
		[]selection.Range{{8, 13}, {6, 15}},
	}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
