package nilslice

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestNilResultIsActuallyNil перевіряє, що NilResult() повертає
// справжній nil-зріз, а не порожній літерал.
func TestNilResultIsActuallyNil(t *testing.T) {
	got := NilResult()
	if got != nil {
		t.Errorf("NilResult() = %#v, want nil", got)
	}
	if len(got) != 0 {
		t.Errorf("len(NilResult()) = %d, want 0", len(got))
	}
}

// TestEmptyResultIsNotNil перевіряє, що EmptyResult() повертає
// порожній, але НЕ nil, зріз.
func TestEmptyResultIsNotNil(t *testing.T) {
	got := EmptyResult()
	if got == nil {
		t.Errorf("EmptyResult() = nil, want a non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("len(EmptyResult()) = %d, want 0", len(got))
	}
}

// TestNilAndEmptyAreNotDeepEqual демонструє те саме, що зазначено
// у промпті домашнього завдання: reflect.DeepEqual вважає
// nil-зріз і порожній зріз РІЗНИМИ значеннями, хоча len() у обох 0.
func TestNilAndEmptyAreNotDeepEqual(t *testing.T) {
	if reflect.DeepEqual(NilResult(), EmptyResult()) {
		t.Error("reflect.DeepEqual(NilResult(), EmptyResult()) = true, want false — " +
			"nil-зріз і порожній зріз не є однаковими для reflect.DeepEqual")
	}
}

// TestJSONMarshalDifference — ключовий тест Завдання 3: показує, що
// відмінність між nil і порожнім зрізом впливає на реальну поведінку —
// маршалинг у JSON.
//
//	nil-зріз       -> "null"
//	порожній зріз  -> "[]"
func TestJSONMarshalDifference(t *testing.T) {
	nilJSON, err := json.Marshal(NilResult())
	if err != nil {
		t.Fatalf("json.Marshal(NilResult()) error: %v", err)
	}
	if string(nilJSON) != "null" {
		t.Errorf("json.Marshal(NilResult()) = %s, want null", nilJSON)
	}

	emptyJSON, err := json.Marshal(EmptyResult())
	if err != nil {
		t.Fatalf("json.Marshal(EmptyResult()) error: %v", err)
	}
	if string(emptyJSON) != "[]" {
		t.Errorf("json.Marshal(EmptyResult()) = %s, want []", emptyJSON)
	}
}

// TestIsNilHelper — базова перевірка допоміжної функції, яку ми
// надали готовою.
func TestIsNilHelper(t *testing.T) {
	if !IsNil(NilResult()) {
		t.Error("IsNil(NilResult()) = false, want true")
	}
	if IsNil(EmptyResult()) {
		t.Error("IsNil(EmptyResult()) = true, want false")
	}
}
