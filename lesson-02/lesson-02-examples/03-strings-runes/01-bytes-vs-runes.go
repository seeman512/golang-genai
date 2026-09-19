// Розділ 2.5 · Рядки, байти й руни
//
// Демонструє: len(s) повертає БАЙТИ, а range повертає РУНИ.
// "Přítel" ("друг" чеською) містить два багатобайтові символи: ř і í.
package main

import "fmt"

func main() {
	s := "Přítel"

	fmt.Println("len(s) in bytes:", len(s)) // 8, а не 6!

	for i, r := range s {
		fmt.Printf("byte-index %d -> rune %q\n", i, r)
	}

	// Перетворюйте явно, коли потрібна одна руна на символ
	runes := []rune(s)
	fmt.Println("rune count:", len(runes)) // 6
}
