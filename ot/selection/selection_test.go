package selection_test

import (
	"encoding/json"
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

func TestMarshal(t *testing.T) {
	for _, tc := range []struct {
		input  *selection.Selection
		output map[string]interface{}
	}{
		{input: &selection.Selection{}, output: map[string]interface{}{"ranges": []map[string]interface{}{}}},
		{input: &selection.Selection{[]selection.Range{}}, output: map[string]interface{}{"ranges": []map[string]interface{}{}}},

		{input: &selection.Selection{[]selection.Range{
			{5, 8},
		}}, output: map[string]interface{}{"ranges": []map[string]interface{}{
			{"anchor": 5, "head": 8},
		}}},

		{input: &selection.Selection{[]selection.Range{
			{2, 11},
			{3, 7},
		}}, output: map[string]interface{}{"ranges": []map[string]interface{}{
			{"anchor": 2, "head": 11},
			{"anchor": 3, "head": 7},
		}}},
	} {
		actual := tc.input.Marshal()
		if expected := tc.output; !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected marshaled %+v to equal %+v, got %+v", tc.input, expected, actual)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	for _, j := range []string{
		`{"foo": "bar"}`,
		`{"ranges": 123}`,
		`{"ranges": "abc"}`,
		`{"ranges": null}`,
		`{"ranges": [{"anchor": 1, "tail": 2}]}`,
		`{"ranges": [{"anchor": 1, "head": "shoulder"}]}`,
	} {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(j), &m)
		if err != nil {
			t.Fatalf("invalid test case")
		}
		_, err = selection.Unmarshal(m)
		if err != selection.ErrUnmarshalFailed {
			t.Errorf("expected ErrUnmarshalFailed when unmarshalling %+v, got %v", m, err)
		}
	}

	for _, tc := range []struct {
		input  string
		output *selection.Selection
	}{
		{input: `{"ranges": []}`, output: &selection.Selection{[]selection.Range{}}},

		{input: `{"ranges": [{"anchor": 1, "head": 3}]}`, output: &selection.Selection{[]selection.Range{
			{1, 3},
		}}},

		{input: `{"ranges": [{"anchor": 2, "head": 5}, {"anchor": 1, "head": 4}]}`, output: &selection.Selection{[]selection.Range{
			{2, 5},
			{1, 4},
		}}},
	} {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(tc.input), &m)
		if err != nil {
			t.Fatalf("invalid test case")
		}

		s, err := selection.Unmarshal(m)
		if err != nil {
			t.Errorf("expected no error, got +%v", err)
		}
		if actual, expected := s, tc.output; !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %+v, got %+v", expected, actual)
		}
	}
}
