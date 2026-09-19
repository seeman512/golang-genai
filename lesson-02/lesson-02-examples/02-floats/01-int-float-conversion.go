// Розділ 2.2 · Числа з рухомою комою
//
// Демонструє: Go НЕ конвертує int в float64 автоматично.
// Розкоментуйте рядок з помилкою, щоб показати помилку компіляції.
package main

import "fmt"

func main() {
	price := 19.99
	quantity := 3

	// total := price * quantity        // ПОМИЛКА КОМПІЛЯЦІЇ: mismatched types
	total := price * float64(quantity) // ОК: явне перетворення

	fmt.Println("total:", total)
}
