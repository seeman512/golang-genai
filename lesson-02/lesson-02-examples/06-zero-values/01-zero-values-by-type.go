// Розділ 4.2 · Таблиця нульових значень
//
// Демонструє: кожен тип автоматично отримує своє нульове значення
// без явної інціалізації — немає стану "не визначено".
package main

import "fmt"

type Point struct{ X, Y int }

func main() {
	var i int
	var f float64
	var s string
	var b bool
	var p *int
	var sl []int
	var m map[string]int
	var pt Point

	fmt.Println("int:", i)
	fmt.Println("float64:", f)
	fmt.Printf("string: %q\n", s)
	fmt.Println("bool:", b)
	fmt.Println("pointer:", p)
	fmt.Println("slice:", sl, "| nil?", sl == nil)
	fmt.Println("map:", m, "| nil?", m == nil)
	fmt.Println("struct:", pt) // {0 0} — кожне поле отримало власне нульове значення
}
