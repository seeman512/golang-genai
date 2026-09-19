// Розділ 3.3 · switch
//
// Демонструє: switch без автоматичного проходження (fall-through)
// та type switch — перемикання за динамічним типом значення.
package main

import (
	"fmt"
	"time"
)

func describe(val interface{}) {
	switch v := val.(type) {
	case int:
		fmt.Println("int:", v)
	case string:
		fmt.Println("string:", v)
	default:
		fmt.Println("unknown type")
	}
}

func main() {
	switch day := time.Now().Weekday(); day {
	case time.Saturday, time.Sunday:
		fmt.Println("weekend")
	default:
		fmt.Println("weekday")
	}

	describe(42)
	describe("hello")
	describe(3.14)
}
