// Розділ 3.2 · for — єдина конструкція циклу
//
// Демонструє: усі чотири форми for, які покривають те, що в інших
// мовах реалізовано окремими ключовими словами (while, do-while, foreach).
package main

import "fmt"

func main() {
	fmt.Println("--- класичний for ---")
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	fmt.Println("--- еквівалент while ---")
	n := 0
	for n < 3 {
		fmt.Println("n =", n)
		n++
	}

	fmt.Println("--- нескінченний цикл (з break) ---")
	count := 0
	for {
		if count >= 3 {
			break
		}
		fmt.Println("count =", count)
		count++
	}

	fmt.Println("--- range-цикл ---")
	items := []string{"apple", "banana", "cherry"}
	for i, v := range items {
		fmt.Println(i, v)
	}
}
