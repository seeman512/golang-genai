// Розділ 2.1 · Цілі числа
//
// Демонструє: типовий int, платформозалежність, та типи з
// фіксованим розміром (int64, uint16) для конкретних задач.
package main

import "fmt"

func main() {
	var age int = 30                   // платформозалежний, типовий вибір
	var fileSize int64 = 5_368_709_120 // 5 ГБ, явно 64-бітний
	var httpStatus uint16 = 404        // явно 16-бітний, не буває від'ємним

	fmt.Println("age:", age)
	fmt.Println("fileSize:", fileSize)
	fmt.Println("httpStatus:", httpStatus)
}
