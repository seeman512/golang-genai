// Розділ 8.4 · Де згенерований ШІ код найчастіше помиляється
//
// Демонструє: типову помилку ШІ-асистента — сеттер з приймачем-
// значенням, який компілюється без жодної помилки, але непомітно
// нічого не робить. Порівняно з правильною версією поруч.
package main

import "fmt"

type Customer struct{ Email string }

// ✗ ТИПОВА ПОМИЛКА ШІ
func (c Customer) SetEmailWrong(e string) {
	c.Email = e // непомітно нічого не робить!
}

// ✓ ПРАВИЛЬНА ВЕРСІЯ
func (c *Customer) SetEmailRight(e string) {
	c.Email = e // змінює реальну структуру
}

func main() {
	cust := Customer{Email: "old@mail.com"}

	cust.SetEmailWrong("new@mail.com")
	fmt.Println("після SetEmailWrong:", cust.Email) // old@mail.com — не змінився!

	cust.SetEmailRight("new@mail.com")
	fmt.Println("після SetEmailRight:", cust.Email) // new@mail.com — оновлено
}
