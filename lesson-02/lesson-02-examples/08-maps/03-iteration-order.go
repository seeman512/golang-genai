// Розділ 6.3 · Пастка порядку ітерації
//
// Демонструє: порядок ітерації мапи рандомізується щоразу під час
// запуску (запустіть цей файл кілька разів і порівняйте перший блок
// виводу) — і як отримати стабільний, детермінований порядок.
package main

import (
	"fmt"
	"sort"
)

func main() {
	prices := map[string]float64{
		"apple":  0.5,
		"banana": 0.3,
		"cherry": 2.0,
	}

	fmt.Println("--- нестабільний порядок (запустіть кілька разів!) ---")
	for key, value := range prices {
		fmt.Println(key, value)
	}

	fmt.Println("--- стабільний, відсортований порядок ---")
	keys := make([]string, 0, len(prices))
	for k := range prices {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, prices[k])
	}
}
