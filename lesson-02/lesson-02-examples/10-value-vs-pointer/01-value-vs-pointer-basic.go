// Розділ 8.1 · Основна відмінність
//
// Демонструє: тип-значення копіює дані при передачі у функцію,
// тип-вказівник дає спільний доступ до оригінальних даних.
package main

import "fmt"

type Customer struct{ Name string }

func updateByValue(c Customer) {
	c.Name = "Changed" // лише локальна копія
}

func updateByPointer(c *Customer) {
	c.Name = "Changed" // змінює оригінал
}

func main() {
	cust := Customer{Name: "Alice"}

	updateByValue(cust)
	fmt.Println("after updateByValue:", cust.Name) // Alice — не змінилося

	updateByPointer(&cust)
	fmt.Println("after updateByPointer:", cust.Name) // Changed
}
