// Розділ 2.5 · Рядки, байти й руни
//
// Демонструє: чому strings.Builder ефективніший за конкатенацію
// через + у циклі (яка щоразу створює новий рядок — O(n²)).
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Погано в циклі: O(n²) — щоразу створюється новий рядок
	var naive string
	for i := 0; i < 5; i++ {
		naive += "x"
	}
	fmt.Println("naive:", naive)

	// Добре: O(n) — Builder переалоковує пам'ять значно рідше
	var b strings.Builder
	for i := 0; i < 5; i++ {
		b.WriteString("x")
	}
	fmt.Println("builder:", b.String())
}
