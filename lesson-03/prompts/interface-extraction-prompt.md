# Промпт для виокремлення інтерфейсу (Розділ 2)

Вставте нижче точний текст промпту, який ви розробили для Розділу 2
домашньої роботи. Промпт має явно містити:

- **роль** ШІ (наприклад, "Ти — досвідчений Go-розробник, що
  спеціалізується на чистій архітектурі...");
- **контекст** — наведений у завданні код `OrderService`;
- **завдання** — виокремити мінімальний інтерфейс `OrderStore`,
  переписати `OrderService` на dependency injection, згенерувати мок;
- **обмеження** — маленький інтерфейс (приказка Роба Пайка),
  конкретний стиль коду, тощо;
- **бажаний формат виводу** (наприклад, "поверни лише Go-код у трьох
  блоках: інтерфейс, сервіс, мок").

## My prompt

You are an experienced Go developer who specializes in clean architecture, consumer-side interfaces, dependency injection, and testable code.

We are refactoring the `orders` package in a small Go project. The current `OrderService` is coupled directly to a database implementation (`*sql.DB`) and its `PlaceOrder` method writes an order using the store's `Exec` operation. Inspect the existing `orders/orders.go` file and preserve its public behavior and the existing method signature:

```go
func (s *OrderService) PlaceOrder(orderID string, amount float64) error
```

Refactor this code as follows:

1. Extract the smallest possible consumer-side interface named `OrderStore`. It must contain only the database operation that `OrderService` actually needs to place an order. Follow Rob Pike's principle: "the bigger the interface, the weaker the abstraction." Do not copy the complete API of `*sql.DB` into the interface.
2. Change `OrderService` to depend on `OrderStore`, not on `*sql.DB` or any other concrete database type.
3. Add a constructor with dependency injection:

   ```go
   func NewOrderService(store OrderStore) *OrderService
   ```

4. Implement or preserve `PlaceOrder` so that it calls `Exec` with an appropriate `INSERT` statement and the `orderID` and `amount` arguments. Return the store error unchanged; return `nil` on success.
5. Generate a lightweight hand-written mock or fake `OrderStore` for unit tests. It must work without a real database, record the query and arguments, count calls, and allow a configured error to be returned. Include enough code to test both a successful order and propagation of a store error.

Constraints:

- Use idiomatic Go and keep the interface minimal.
- Keep the service focused on order logic; do not add database connection, global state, reflection, or unrelated abstractions.
- Use the repository's existing `Exec` signature and Go version. Do not add or update dependencies.
- Do not change production behavior beyond replacing the concrete database dependency with the interface.
- Preserve the package name and existing public API unless a change is required by the dependency-injection design.
- Do not modify existing tests; the resulting code must be compatible with them.
- Explain any assumption if the current starter code differs from the context above.

Return the result in exactly three clearly labeled Go code blocks, in this order:

1. `OrderStore interface` — the minimal interface.
2. `OrderService` — the service, constructor, and `PlaceOrder` implementation.
3. `Mock/fake and tests` — the test double and focused unit tests.

Before returning the code, verify that it compiles, that the mock satisfies `OrderStore`, and that errors from `Exec` are propagated correctly. Do not include production dependencies on a real database in the mock or tests.

