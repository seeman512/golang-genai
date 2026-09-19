// Розділ 8.2 · Приймачі методів: значення проти вказівника
//
// Демонструє: приймач-значення працює з копією (зміна втрачається),
// приймач-вказівник змінює оригінал. Зверніть увагу: Go автоматично
// бере &c, коли викликає метод із приймачем-вказівником на змінній.
package main

import "fmt"

type Counter struct{ n int }

// Приймач-значення: працює з КОПІЄЮ — Counter викликача не змінюється
func (c Counter) IncrementCopy() { c.n++ }

// Приймач-вказівник: працює з ОРИГІНАЛОМ через його адресу
func (c *Counter) Increment() { c.n++ }

func main() {
	c := Counter{}

	c.IncrementCopy()
	fmt.Println("after IncrementCopy:", c.n) // 0 — копію було відкинуто

	c.Increment()                        // Go автоматично бере &c тут
	fmt.Println("after Increment:", c.n) // 1
}
