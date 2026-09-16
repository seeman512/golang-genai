// Команда greet зчитує ім'я користувача з аргументів командного рядка
// і виводить привітання, згенероване функцією greet.Greet.
package main

import (
	"fmt"
	"os"

	"lesson01/greet"
)

func main() {
	name := ""
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	fmt.Println(greet.Greet(name))
}
