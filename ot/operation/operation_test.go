package operation_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nitrous-io/ot.go/ot/operation"
)

func TestNew(t *testing.T) {
	top := operation.New()

	if actual := reflect.TypeOf(top); actual != reflect.TypeOf(&operation.Operation{}) {
		t.Fatalf("expected New to return a pointer to operation.Operation, got %v", actual)
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
	top := &operation.Operation{}

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

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{N: 2},
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

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{N: 6},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestDelete(t *testing.T) {
	top := &operation.Operation{}

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

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{N: -2},
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

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{N: -5},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestInsert(t *testing.T) {
	top := &operation.Operation{}

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

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{S: "foo"},
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

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{S: "foobarbaz"},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestMultipleOps(t *testing.T) {
	top := &operation.Operation{}

	top.Retain(1).Delete(2).Delete(1).Retain(2).Retain(3).Delete(1)

	if actual, expected := top.BaseLen, 10; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 6; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{N: 1},
		&operation.Op{N: -3},
		&operation.Op{N: 5},
		&operation.Op{N: -1},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	top = &operation.Operation{}

	top.Retain(1).Insert("foo").Delete(2).Insert("bar").Delete(1).Retain(2).Retain(3).Insert("baz").Delete(1).Delete(3)

	if actual, expected := top.BaseLen, 13; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 15; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{N: 1},
		&operation.Op{S: "foobar"},
		&operation.Op{N: -3},
		&operation.Op{N: 5},
		&operation.Op{S: "baz"},
		&operation.Op{N: -4},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestApply(t *testing.T) {
	top := &operation.Operation{}
	top.Retain(3)

	_, err := top.Apply("fo")

	if err != operation.ErrBaseLenMismatch {
		t.Errorf("expected operation.ErrBaseLenMismatch, got %v", err)
	}

	_, err = top.Apply("food")

	if err != operation.ErrBaseLenMismatch {
		t.Errorf("expected operation.ErrBaseLenMismatch, got %v", err)
	}

	s, err := top.Apply("foo")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "foo"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	top = &operation.Operation{}
	s, err = top.Insert("bar").Apply("")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "bar"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	top = &operation.Operation{}
	s, err = top.Delete(3).Apply("baz")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual := s; actual != "" {
		t.Errorf("expected empty string, got %s", actual)
	}

	top = &operation.Operation{}
	s, err = top.Retain(1).Insert("dar").Delete(1).Retain(1).Insert("biz").Apply("fox")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "fdarxbiz"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestTransform(t *testing.T) {
	a := &operation.Operation{}
	a.Retain(1)

	b := &operation.Operation{}
	b.Retain(2)

	_, _, err := operation.Transform(a, b)

	if err != operation.ErrBaseLenMismatch {
		t.Errorf("expected ErrBaseLenMismatch, got %v", err)
	}

	// apply(apply(S, A), B') = apply(apply(S, B), A')

	s := "She is a girl!!!"
	o := "He was one beautiful man."

	a = &operation.Operation{}
	a.Retain(4).Delete(1).Insert("wa").Retain(4).Insert("beautiful ").Retain(4).Delete(3).Insert(".")

	b = &operation.Operation{}
	b.Delete(2).Insert("H").Retain(5).Delete(1).Insert("one").Retain(1).Delete(4).Insert("man").Delete(2).Retain(1)

	a1, b1, err := operation.Transform(a, b)

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

func TestMarshal(t *testing.T) {
	top := &operation.Operation{}
	top.Retain(2).Insert("H").Retain(5).Insert("one").Delete(1).Retain(1).Insert("man").Delete(6).Retain(1)

	ops := top.Marshal()

	if actual, expected := ops, []interface{}{2, "H", 5, "one", -1, 1, "man", -6, 1}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestUnmarshal(t *testing.T) {
	j := `[1, ["H"], -1]`
	var ops []interface{}
	err := json.Unmarshal([]byte(j), &ops)
	if err != nil {
		t.Fatalf("test case error")
	}

	_, err = operation.Unmarshal(ops)
	if err != operation.ErrUnmarshalFailed {
		t.Fatalf("expected ErrUnmarshalFailed, got %v", err)
	}

	j = `[2, "Sh", 5, -1, "one", 1, -4, "man", -2, 1]`
	err = json.Unmarshal([]byte(j), &ops)
	if err != nil {
		t.Fatalf("test case error")
	}

	top, err := operation.Unmarshal(ops)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if top.Ops == nil {
		t.Fatalf("expected list of ops not to be nil, got nil")
	}

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{N: 2},
		&operation.Op{S: "Sh"},
		&operation.Op{N: 5},
		&operation.Op{S: "one"},
		&operation.Op{N: -1},
		&operation.Op{N: 1},
		&operation.Op{S: "man"},
		&operation.Op{N: -6},
		&operation.Op{N: 1},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	if actual, expected := top.BaseLen, 16; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 17; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}
}
