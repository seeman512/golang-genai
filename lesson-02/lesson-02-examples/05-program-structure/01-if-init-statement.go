// Розділ 3.1 · if / else
//
// Демонструє: форму "if з інструкцією-ініціалізатором" — ідіоматичний
// Go-спосіб перевіряти помилку одразу після її отримання, звужуючи
// область видимості v і err лише до цього блоку.
package main

import (
	"errors"
	"fmt"
	"log"
)

func computeValue() (int, error) {
	return 150, nil
}

func computeValueWithError() (int, error) {
	return 0, errors.New("щось пішло не так")
}

func main() {
	if v, err := computeValue(); err != nil {
		log.Fatal(err)
	} else if v > 100 {
		fmt.Println("large value:", v)
	}

	if _, err := computeValueWithError(); err != nil {
		fmt.Println("отримали очікувану помилку:", err)
	}
}
