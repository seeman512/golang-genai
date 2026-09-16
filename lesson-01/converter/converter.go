// Package converter реалізує просту конвертацію валют за фіксованим курсом.
//
// Завдання 4 (бонус): реалізуйте ConvertCurrency самостійно або з допомогою
// ШІ-асистента — і порівняйте результат із вимогами в converter_test.go.
package converter

import "errors"

// ConvertCurrency конвертує amount за курсом rate.
//
// Правила:
//   - результат дорівнює amount * rate;
//   - якщо amount від'ємний — повертається помилка;
//   - якщо rate від'ємний або дорівнює нулю — повертається помилка.
func ConvertCurrency(amount float64, rate float64) (float64, error) {
	if amount < 0 {
		return 0, errors.New("amount must not be negative")
	}
	if rate <= 0 {
		return 0, errors.New("rate must be positive")
	}

	return amount * rate, nil
}
