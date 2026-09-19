// Розділ 6.2 · Нульове значення мапи
//
// Демонструє: читання з nil-мапи безпечне, запис — панікує.
// Використано recover(), щоб показати паніку без аварійного завершення.
package main

import "fmt"

func main() {
	var m map[string]int                       // nil
	fmt.Println("читання з nil-мапи:", m["x"]) // 0 — безпечно!

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("сталася паніка:", r)
			}
		}()
		m["x"] = 1 // ПАНІКА: запис у nil-мапу
	}()
}
