// Package deadlock містить практичне застосування аналізу з
// Частини 2 домашньої роботи: той самий дедлок із завдання,
// але як тестована функція, яку потрібно виправити.
//
// Оригінальний баг (для довідки — саме цей код ви аналізуєте у
// deadlock-analysis.md):
//
//	var wg sync.WaitGroup
//	ch := make(chan int)
//	wg.Add(1)
//	go func() {
//	    ch <- 42
//	    wg.Done()
//	}()
//	wg.Wait()       // <-- блокується тут: горутина застрягла на
//	                //     "ch <- 42", бо ніхто ще не читає з ch
//	fmt.Println(<-ch)
package deadlock

// Run returns 42 sent by a goroutine through a channel without deadlocking.
// The receive operation synchronizes with the send, so a WaitGroup is not
// needed when there is only one value to exchange.
func Run() int {
	ch := make(chan int)

	go func() {
		ch <- 42
	}()

	return <-ch
}
