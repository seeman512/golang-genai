// Розділ 4.4 · Пастки nil
//
// Демонструє: на відміну від мапи, append() для nil-зрізу працює
// НОРМАЛЬНО і виділяє пам'ять за потреби. Паніки не буде.
package main

import "fmt"

func main() {
	var s []int
	fmt.Println("nil-зріз перед append:", s, "| nil?", s == nil)

	s = append(s, 1) // OK — жодної паніки
	fmt.Println("після append:", s, "| nil?", s == nil)
}
