// Розділ 2.6 · Константи та iota
//
// Демонструє: ідіоматичний спосіб реалізувати enum у Go — через iota
// плюс метод String(), щоб fmt.Println друкував назву, а не число.
package main

import "fmt"

type Weekday int

const (
	Sunday    Weekday = iota // 0
	Monday                   // 1
	Tuesday                  // 2
	Wednesday                // 3
	Thursday                 // 4
	Friday                   // 5
	Saturday                 // 6
)

func (d Weekday) String() string {
	return [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}[d]
}

func main() {
	today := Wednesday
	fmt.Println("today:", today)       // Wed — завдяки String()
	fmt.Println("as int:", int(today)) // 3
}
