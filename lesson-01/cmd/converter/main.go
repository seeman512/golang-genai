// Команда converter зчитує суму та курс з аргументів командного рядка
// і виводить результат конвертації, застосовуючи converter.ConvertCurrency.
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"example.com/lesson01/converter"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: converter <amount> <rate>")
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		log.Fatalf("invalid amount: %v", err)
	}

	rate, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatalf("invalid rate: %v", err)
	}

	result, err := converter.ConvertCurrency(amount, rate)
	if err != nil {
		log.Fatalf("conversion error: %v", err)
	}

	fmt.Printf("%.2f = %.2f\n", amount, result)
}
