// Розділ 6.1 · Що таке мапа
//
// Демонструє: базові операції з мапою та ідіому "comma ok" для
// перевірки наявності ключа.
package main

import "fmt"

func main() {
	prices := map[string]float64{
		"apple":  0.5,
		"banana": 0.3,
	}

	prices["cherry"] = 2.0 // вставка/оновлення

	v, ok := prices["kiwi"] // ok == false, v == 0 (нульове значення)
	fmt.Println("kiwi:", v, "| присутній?", ok)

	delete(prices, "banana")
	fmt.Println("після delete:", prices)

	inventory := map[string]int{"apples": 10}
	qty, ok2 := inventory["kiwi"]
	if !ok2 {
		fmt.Println("kiwi не знайдено в інвентарі")
	}
	fmt.Println("qty:", qty) // 0, навіть якщо ok2 == false
}
