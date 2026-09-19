// Розділ 5.2 · Зрізи: анатомія та семантика посилань
//
// Демонструє: копіювання зрізу копіює лише ЗАГОЛОВОК {ptr,len,cap},
// а не дані. Обидва зрізи бачать той самий базовий масив, доки один
// з них не переросте свою місткість через append.
package main

import "fmt"

func main() {
	scores := []int{90, 85, 80}
	other := scores // копіюється лише ЗАГОЛОВОК — дані спільні!
	other[0] = 100

	fmt.Println("scores:", scores) // [100 85 80] — теж змінилося!
	fmt.Println("other: ", other)

	scores = append(scores, 60) // може вирости й отримати НОВИЙ базовий масив
	fmt.Println("after append, scores:", scores)
	fmt.Println("other залишився без змін:", other)
}
