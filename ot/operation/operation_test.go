package operation_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nitrous-io/ot.go/ot"
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
	ot.TextEncoding = ot.TextEncodingTypeUTF8
	defer func() {
		ot.TextEncoding = ot.TextEncodingTypeUTF8
	}()

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
		&operation.Op{S: []rune("foo")},
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

	top.Insert("유니코드😄")

	if actual, expected := top.BaseLen, 0; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 14; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{S: []rune("foobarbaz유니코드😄")},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	// utf-16
	ot.TextEncoding = ot.TextEncodingTypeUTF16
	top = operation.New().Insert("abc").Insert("가나다").Insert("αβγ").Insert("😄😃😀")

	if actual, expected := top.BaseLen, 0; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 15; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{S: []rune{97, 98, 99, 44032, 45208, 45796, 945, 946, 947, 55357, 56836, 55357, 56835, 55357, 56832}},
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
		&operation.Op{S: []rune("foobar")},
		&operation.Op{N: -3},
		&operation.Op{N: 5},
		&operation.Op{S: []rune("baz")},
		&operation.Op{N: -4},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestApply(t *testing.T) {
	ot.TextEncoding = ot.TextEncodingTypeUTF8
	defer func() {
		ot.TextEncoding = ot.TextEncodingTypeUTF8
	}()

	// retain

	top := operation.New().Retain(3)

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

	s, err = top.Apply("사랑💖")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "사랑💖"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	// insert

	s, err = operation.New().Insert("bar").Apply("")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "bar"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	s, err = operation.New().Insert("愛してる💖").Apply("")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "愛してる💖"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	// delete

	s, err = operation.New().Delete(3).Apply("baz")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual := s; actual != "" {
		t.Errorf("expected empty string, got %s", actual)
	}

	s, err = operation.New().Delete(2).Retain(3).Apply("🌍你好世界")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "好世界"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	// combined

	s, err = operation.New().Retain(2).Insert("far").Delete(1).Retain(1).Insert("대성공💯").Apply("🐺dog")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "🐺dfarg대성공💯"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}

	// utf-16
	ot.TextEncoding = ot.TextEncodingTypeUTF16

	s, err = operation.New().Retain(3).Insert("far").Delete(1).Retain(1).Insert("대성공💯").Apply("🐺dog")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if actual, expected := s, "🐺dfarg대성공💯"; actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestTransform(t *testing.T) {
	ot.TextEncoding = ot.TextEncodingTypeUTF8
	defer func() {
		ot.TextEncoding = ot.TextEncodingTypeUTF8
	}()

	a := operation.New().Retain(1)

	b := operation.New().Retain(2)

	_, _, err := operation.Transform(a, b)

	if err != operation.ErrBaseLenMismatch {
		t.Errorf("expected ErrBaseLenMismatch, got %v", err)
	}

	// apply(apply(S, A), B') = apply(apply(S, B), A')

	testTransform := func(s, o string, a, b *operation.Operation) {
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

	s := "She is a girl!!!"
	o := "He was 정말로 beautiful man."
	a = operation.New().Retain(4).Delete(1).Insert("wa").Retain(4).Insert("beautiful ").Retain(4).Delete(3).Insert(".")
	b = operation.New().Delete(2).Insert("H").Retain(5).Delete(1).Insert("정말로").Retain(1).Delete(4).Insert("man").Delete(2).Retain(1)

	testTransform(s, o, a, b)

	// utf-16
	ot.TextEncoding = ot.TextEncodingTypeUTF16

	s = "She is 😝 girl👧!"
	o = "He was 👍👍 beautiful man."
	a = operation.New().Retain(4).Delete(1).Insert("wa").Retain(5).Insert("beautiful ").Retain(4).Delete(3).Insert(".")
	b = operation.New().Delete(2).Insert("H").Retain(5).Delete(2).Insert("👍👍").Retain(1).Delete(4).Insert("man").Delete(2).Retain(1)

	testTransform(s, o, a, b)
}

func TestMarshal(t *testing.T) {
	ot.TextEncoding = ot.TextEncodingTypeUTF8
	defer func() {
		ot.TextEncoding = ot.TextEncodingTypeUTF8
	}()

	top := operation.New().Retain(2).Insert("H").Retain(5).Insert("정말로").Delete(1).Retain(1).Insert("man").Delete(6).Retain(1)
	ops := top.Marshal()

	if actual, expected := ops, []interface{}{2, "H", 5, "정말로", -1, 1, "man", -6, 1}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	// utf-16
	ot.TextEncoding = ot.TextEncodingTypeUTF16

	top = operation.New().Insert("abc").Retain(1).Insert("가나다").Retain(2).Insert("αβγ").Retain(3).Insert("😄😃😀")
	ops = top.Marshal()

	if actual, expected := ops, []interface{}{"abc", 1, "가나다", 2, "αβγ", 3, "😄😃😀"}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestUnmarshal(t *testing.T) {
	ot.TextEncoding = ot.TextEncodingTypeUTF8
	defer func() {
		ot.TextEncoding = ot.TextEncodingTypeUTF8
	}()

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

	j = `[2, "Sh", 5, -1, "정말😄", 1, -4, "man", -2, 1]`
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
		&operation.Op{S: []rune("Sh")},
		&operation.Op{N: 5},
		&operation.Op{S: []rune("정말😄")},
		&operation.Op{N: -1},
		&operation.Op{N: 1},
		&operation.Op{S: []rune("man")},
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

	// utf-16
	ot.TextEncoding = ot.TextEncodingTypeUTF16

	j = `["abc", 1, "가나다", 2, "αβγ", 3, "😄😃😀"]`
	err = json.Unmarshal([]byte(j), &ops)
	if err != nil {
		t.Fatalf("test case error")
	}

	top, err = operation.Unmarshal(ops)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if top.Ops == nil {
		t.Fatalf("expected list of ops not to be nil, got nil")
	}

	if actual, expected := top.Ops, []*operation.Op{
		&operation.Op{S: []rune{97, 98, 99}},
		&operation.Op{N: 1},
		&operation.Op{S: []rune{44032, 45208, 45796}},
		&operation.Op{N: 2},
		&operation.Op{S: []rune{945, 946, 947}},
		&operation.Op{N: 3},
		&operation.Op{S: []rune{55357, 56836, 55357, 56835, 55357, 56832}},
	}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}

	if actual, expected := top.BaseLen, 6; actual != expected {
		t.Errorf("expected base length of %d, got %d", expected, actual)
	}

	if actual, expected := top.TargetLen, 21; actual != expected {
		t.Errorf("expected target length of %d, got %d", expected, actual)
	}
}
