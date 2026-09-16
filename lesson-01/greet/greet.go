// Package greet містить логіку привітання користувача.
//
// Завдання 2: реалізуйте функцію Greet самостійно, вручну, без ШІ.
// Автоматичні тести дивіться у файлі greet_test.go.
package greet

import "strings"

// Greet приймає ім'я користувача і повертає рядок привітання.
//
// Правила:
//   - для непорожнього імені (пробіли на краях обрізаються):
//     "Hello, <name>! Welcome to Go."
//   - якщо після обрізання пробілів ім'я порожнє:
//     "Hello, stranger! Welcome to Go."
func Greet(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "stranger"
	}
	return "Hello, " + name + "! Welcome to Go."
}
