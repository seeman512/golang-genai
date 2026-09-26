# Аналіз дедлоку (Частина 2)

## Код, що аналізується

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan int) // небуферизований канал

	wg.Add(1)
	go func() {
		// симуляція важкої роботи
		ch <- 42
		wg.Done()
	}()

	wg.Wait()
	fmt.Println(<-ch)
}
```

## Мій промпт до ШІ

Вставте нижче точний текст промпту, який ви розробили. Промпт має
явно вимагати від ШІ:

- **не просто переписати код**, а пояснити точну послідовність
  переплетення (interleaving) горутин, що веде до зависання;
- проаналізувати, чому виклик `wg.Wait()` перед читанням з
  небуферизованого каналу зупиняє виконання.

### Мій промпт

You are an expert in Go concurrency. Analyze the issue in
`lesson-04/deadlock/fixed.go` and the related tests in
`lesson-04/deadlock/fixed_test.go`. Pay special attention to the original
code using an unbuffered channel and `sync.WaitGroup`: the goroutine tries
to execute `ch <- 42`, while the main goroutine calls `wg.Wait()` first and
only reads from the channel afterward. Also explain why the current `Run`
stub, which creates a channel and immediately executes `<-ch` without a
sender goroutine, also deadlocks.

Do not simply rewrite the code. First, describe step by step the exact
interleaving of the goroutines and every operation at which they become
blocked. Explain the role of the unbuffered channel, the point at which the
`WaitGroup` blocks, and why neither goroutine can reach the operation that
would unblock the other.

Then propose and implement the minimal fix specifically in the `Run`
function. The fix must start a sender goroutine, receive the value `42`,
and finish without a deadlock. Do not change the `Run` signature, modify
the tests, or add global shared state. Explain why a `WaitGroup` is
unnecessary when the function exchanges only one value and why the direct
channel receive prevents the deadlock. Run `gofmt` and
`go test -race ./deadlock/...` to verify the fix, and present the problem
explanation, the changed code fragment, and the verification result as
separate parts of your response.

## Пояснення від ШІ

Вставте нижче відповідь ШІ-асистента на ваш промпт — покрокове
пояснення переплетення горутин.

### Відповідь ШІ

The original program deadlocks through a circular wait between the main
goroutine and the sender goroutine:

1. The `WaitGroup` counter is set to `1`, and the sender goroutine starts.
2. The sender reaches `ch <- 42`. Because `ch` is unbuffered, this send
   can complete only when another goroutine is ready to receive the value.
   The sender therefore blocks at the send operation.
3. The main goroutine calls `wg.Wait()`. The counter is still `1` because
   the sender has not completed the send and has not reached `wg.Done()`.
   The main goroutine therefore blocks in `Wait`.
4. The receive operation `fmt.Println(<-ch)` is after `wg.Wait()`, so the
   main goroutine cannot reach it. The sender cannot reach `wg.Done()` until
   its send completes, and the send cannot complete until the main goroutine
   receives. Each goroutine is waiting for the other, so the program never
   makes progress.

The original `Run` stub has an even simpler deadlock: it creates an
unbuffered channel and immediately waits to receive from it, but it never
starts a goroutine that sends a value. Consequently, there is no possible
sender to match the receive operation.

The simplified fix removes the `WaitGroup` entirely because there is only
one value to exchange. It starts a sender goroutine and immediately performs
the receive in `Run`. The receive synchronizes with the send on the
unbuffered channel, allowing the sender goroutine to finish; therefore, no
additional wait operation is necessary and `Run` returns `42` without
leaving a blocked goroutine.

## Мінімальне виправлення

Робочу, тестовану версію виправлення реалізовано в
[`deadlock/fixed.go`](./deadlock/fixed.go) (функція `Run`), а
автоматична перевірка — у `deadlock/fixed_test.go`. Коротко опишіть
своїми словами, яку зміну порядку операцій (чи структури каналу) ви
застосували і чому саме вона усуває циклічне очікування.

`Run` now starts a goroutine that sends `42` and immediately receives the
value from the unbuffered channel. The send and receive rendezvous directly,
so the sender can finish and no `sync.WaitGroup` is needed for this single
value exchange. This removes the circular wait and lets the function finish
safely.
