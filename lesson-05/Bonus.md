# Bonus: Error-handling review of Lessons 1–4

## Scope and summary

I reviewed the Go source in `lesson-01/` through `lesson-04/` for ignored errors, lost error context, and validation implemented with `panic`.

- There is **no `_ = err` and no raw `panic(...)` used for validation** in the production Go source of Lessons 1–4. The registration validator in `lesson-05/task1_validation/validator.go` also returns a `ValidationError`; it does not panic.
- One test deliberately discards an error result: `lesson-03/orders/orders_test.go`, in `TestOrderServiceUsesInjectedStore`.
- Two production boundaries pass dependency errors through without adding operation context: `OrderService.PlaceOrder` and `ProcessPayment`. Passing an error through unchanged is not the same as ignoring it, and can be the right choice when the API contract calls for exact propagation. If callers benefit from knowing which operation failed, wrap it with `%w`.
- The other reviewed code either has no error-returning operations or handles/reports the errors it receives. In particular, `log.Fatalf("...: %v", err)` in the command-line converter is a top-level report-and-exit, not a returned error that needs `%w`.

## Findings and proposed fixes

### 1. Ignored error result in an order-service test

**Location:** `lesson-03/orders/orders_test.go`, `TestOrderServiceUsesInjectedStore`

```go
_ = svc.PlaceOrder("order-3", 5)
```

This test is checking that the injected store is called, but it silently discards a possible failure from `PlaceOrder`. A failing call could go unnoticed, making the test pass for the wrong reason.

**Fix:** check the result while retaining the dependency-injection assertion:

```go
if err := svc.PlaceOrder("order-3", 5); err != nil {
    t.Fatalf("PlaceOrder returned an unexpected error: %v", err)
}
if store.calls != 1 {
    t.Errorf("Exec called %d times, want 1", store.calls)
}
```

Discard an error only when it is genuinely irrelevant and that choice is intentional. For tests, explicitly assert either success or the expected error; do not use a blank identifier as a substitute for deciding what the test should verify.

### 2. Store error returned without operation context

**Location:** `lesson-03/orders/orders.go`, `(*OrderService).PlaceOrder`

The method returns the `Exec` error directly. This does preserve the original error, and the lesson's prompt specifically asks for unchanged propagation. If the service is also responsible for identifying the failed operation, wrap the store error so callers retain both context and error identity:

```go
func (s *OrderService) PlaceOrder(orderID string, amount float64) error {
    if err := s.store.Exec(
        "INSERT INTO orders (id, amount) VALUES (?, ?)",
        orderID,
        amount,
    ); err != nil {
        return fmt.Errorf("orders: place order %q: %w", orderID, err)
    }
    return nil
}
```

Add `fmt` to the imports. `%w` is important: callers can still use `errors.Is` or `errors.As` on the underlying store error. A focused test should use a sentinel error from the fake store and assert `errors.Is(err, sentinel)`. If exact, unadorned propagation is an explicit contract and the service adds no useful context, returning `err` unchanged is also valid; wrapping is not mandatory at every layer.

### 3. Payment error returned without operation context

**Location:** `lesson-03/payment/payment.go`, `ProcessPayment`

`method.Pay` errors are propagated rather than ignored, but the returned error does not identify that payment processing failed. If this boundary should add that context, use:

```go
func ProcessPayment(method PaymentMethod, amount float64) error {
    if err := method.Pay(amount); err != nil {
        return fmt.Errorf("payment: process transaction: %w", err)
    }
    fmt.Printf("Payment successful. %s\n", method.LogInfo())
    return nil
}
```

`fmt` is already imported in this file. Preserve the cause with `%w`, not `%v`, so callers can classify it with `errors.Is`/`errors.As`. As with `PlaceOrder`, unchanged propagation is acceptable if callers already have enough context and the function's contract favors direct propagation.

## Validation: return errors, do not panic

No validation `panic` was found in Lessons 1–4. For invalid user or function input, return an error (or a typed validation error) because invalid input is an expected condition that the caller should be able to handle. `panic` is for unrecoverable programmer/runtime invariants, not routine validation.

For example, a validation function can return an explanatory error:

```go
func ValidateAge(age int) error {
    if age < 0 || age > 150 {
        return fmt.Errorf("validate age: %d is outside the allowed range [0, 150]", age)
    }
    return nil
}
```

If callers need to classify this failure, define a sentinel and wrap it with `%w`, or return a typed error. The existing `lesson-05/task1_validation` approach—returning `*ValidationError` with the invalid field names—is suitable for collecting multiple field failures. Do not replace that return path with `panic`.

## Error-handling rules for future changes

1. Check every returned error. Handle it, return it, or deliberately document why it is safe to ignore; avoid `_ = err` as a general workaround.
2. Add context at meaningful abstraction boundaries with `fmt.Errorf("operation: details: %w", err)`. Use `%w` when the caller should still be able to inspect the cause; use `%v` for formatting/logging only when preserving the error chain is not needed.
3. Avoid wrapping repeatedly with identical or unhelpful messages. A wrapper should add information such as the operation, resource, identifier, or attempt.
4. Return validation failures as errors or typed errors. Reserve `panic` for truly unrecoverable invariant violations.
5. At an executable boundary, report the error and exit/return. At a library boundary, generally return the error to the caller instead of logging and continuing.
6. Test both the human-readable context and error identity when wrapping matters (`errors.Is`/`errors.As`), as well as successful behavior.
