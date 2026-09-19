// Розділ 2.6 · Константи та iota
//
// Демонструє: нетипізовані константи адаптуються до контексту —
// нетипізована ціла const спокійно підходить туди, де очікується float64.
package main

import "fmt"

const factor = 2 // нетипізована константа

func main() {
	var price float64 = 19.99
	result := price * factor // OK — factor адаптувався до float64
	fmt.Println("result:", result)
}
