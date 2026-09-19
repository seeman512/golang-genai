// Розділ 10 · Самоперевірка, практичне завдання
//
// Демонструє: розв'язок вправи "напишіть Median(nums []float64)",
// яка коректно обробляє порожній зріз, один елемент та парну кількість.
package main

import (
	"errors"
	"fmt"
	"sort"
)

func Median(nums []float64) (float64, error) {
	if len(nums) == 0 {
		return 0, errors.New("median of empty slice is undefined")
	}
	sorted := make([]float64, len(nums))
	copy(sorted, nums)
	sort.Float64s(sorted)

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2, nil
	}
	return sorted[mid], nil
}

func main() {
	fmt.Println(Median([]float64{3, 1, 2}))    // 2 <nil>
	fmt.Println(Median([]float64{4, 1, 3, 2})) // 2.5 <nil>
	fmt.Println(Median([]float64{5}))          // 5 <nil>
	fmt.Println(Median([]float64{}))           // 0 median of empty slice is undefined
}
