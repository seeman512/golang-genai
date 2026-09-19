// Розділ 5.3 · make, append і зростання місткості
//
// Демонструє: як cap() зростає (часто вдвічі) щоразу, коли len
// перевищує поточну місткість базового масиву.
package main

import "fmt"

func main() {
	s := make([]int, 0, 4) // len=0, cap=4 — місце заброньоване наперед
	fmt.Println("start:", len(s), cap(s))

	for i := 0; i < 6; i++ {
		s = append(s, i)
		fmt.Println("len:", len(s), "cap:", cap(s))
	}
}
