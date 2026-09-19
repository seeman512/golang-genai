// Розділ 7.4 · Вбудовування проти іменованого вкладеного поля
//
// Демонструє: вбудована структура (Person у Employee) просуває поля
// нагору автоматично; іменоване поле (Customer в Order) вимагає
// явного шляху доступу.
package main

import "fmt"

type Person struct {
	Name   string
	Salary int
}

// Вбудована структура: анонімне поле, лише тип без імені
type Employee struct {
	Person // вбудовано — поля Person просунуто нагору
	Title  string
}

type Customer struct {
	Name string
}

// Іменоване вкладене поле
type Order struct {
	Customer Customer // іменоване поле, НЕ вбудоване
}

func main() {
	e := Employee{Person: Person{Name: "Ada", Salary: 90000}, Title: "Engineer"}
	fmt.Println("вбудоване поле напряму:", e.Name) // працює без e.Person.Name!

	order := Order{Customer: Customer{Name: "Bob"}}
	fmt.Println("іменоване поле — потрібен повний шлях:", order.Customer.Name)
}
