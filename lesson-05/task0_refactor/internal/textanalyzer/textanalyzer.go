// Package textanalyzer реалізує прості функції аналізу тексту для
// Завдання 0 (об'єднання text analyzer із Заняття 2 в multi-package проєкт).
package textanalyzer

import (
	"errors"
	"fmt"
	"strings"
)

// ErrEmptyText — sentinel error для порожнього вхідного тексту.
var ErrEmptyText = errors.New("textanalyzer: empty text")

// WordCount повертає кількість слів у text (розділених пробільними
// символами). Якщо text порожній (або складається лише з пробілів),
// повертає помилку, що через %w обгортає ErrEmptyText.
//
// (Завдання 0): реалізуйте.
func WordCount(text string) (int, error) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return 0, fmt.Errorf("textanalyzer: word count: %w", ErrEmptyText)
	}

	return len(words), nil
}

// CharCount повертає кількість символів (рун) у text, без урахування
// пробільних символів на початку/в кінці.
//
// (Завдання 0): реалізуйте.
func CharCount(text string) int {
	trimmed := strings.TrimSpace(text)

	return len([]rune(trimmed))
}
