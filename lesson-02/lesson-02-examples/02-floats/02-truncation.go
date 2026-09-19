// Розділ 2.2 · Числа з рухомою комою
//
// Демонструє: перетворення float -> int ВІДКИДАЄ дробову частину
// в бік нуля, воно НЕ округляє. Для округлення потрібен math.Round().
package main

import (
	"fmt"
	"math"
)

func main() {
	var f float64 = 3.99
	var i int = int(f)
	fmt.Println("int(3.99) =", i) // 3, а не 4!

	rounded := math.Round(f)
	fmt.Println("math.Round(3.99) =", rounded) // 4
}
