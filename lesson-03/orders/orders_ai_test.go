package orders

import (
	"reflect"
	"testing"
)

// mockOrderStore is a test double for OrderStore. It records the Exec call so
// PlaceOrder can be tested without connecting to a real database.
type mockOrderStore struct {
	execCalls int
	query     string
	args      []any
	err       error
}

func (m *mockOrderStore) Exec(query string, args ...any) error {
	m.execCalls++
	m.query = query
	m.args = append([]any(nil), args...)
	return m.err
}

func TestPlaceOrderAI_Success(t *testing.T) {
	store := &mockOrderStore{}
	service := NewOrderService(store)

	if err := service.PlaceOrder("order-ai-1", 125.75); err != nil {
		t.Fatalf("PlaceOrder returned an unexpected error: %v", err)
	}

	if store.execCalls != 1 {
		t.Fatalf("Exec was called %d times, want 1", store.execCalls)
	}

	const wantQuery = "INSERT INTO orders (id, amount) VALUES (?, ?)"
	if store.query != wantQuery {
		t.Fatalf("Exec query = %q, want %q", store.query, wantQuery)
	}

	wantArgs := []any{"order-ai-1", 125.75}
	if !reflect.DeepEqual(store.args, wantArgs) {
		t.Fatalf("Exec args = %#v, want %#v", store.args, wantArgs)
	}
}
