// Розділ 7.2 · Моделювання домену: вкладені структури
//
// Демонструє: композицію "has-a" (Order містить Customer) і зріз
// вкладених структур (Order.Items []Product).
package main

import (
	"fmt"
	"time"
)

type Address struct {
	City, Street string
	ZIP          string
}

type Customer struct {
	ID      int
	Name    string
	Address Address // іменована вкладена структура
}

type Product struct {
	SKU   string
	Name  string
	Price float64
}

type Order struct {
	ID        int
	Customer  Customer  // композиція: Order «has-a» Customer
	Items     []Product // зріз вкладених структур
	CreatedAt time.Time
}

func main() {
	order := Order{
		ID: 1001,
		Customer: Customer{
			ID:   1,
			Name: "Alice",
			Address: Address{
				City:   "Kyiv",
				Street: "Khreshchatyk 1",
				ZIP:    "01001",
			},
		},
		Items: []Product{
			{SKU: "A1", Name: "Mug", Price: 9.99},
			{SKU: "B2", Name: "Notebook", Price: 4.50},
		},
		CreatedAt: time.Now(),
	}

	fmt.Printf("%+v\n", order)
	fmt.Println("customer city:", order.Customer.Address.City)
}
