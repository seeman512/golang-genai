// Package library моделює простий домен «Бібліотека».
//
// Завдання 1: спроєктуйте структури Author, Book, Library самостійно
// (або з допомогою ШІ) і реалізуйте метод AddBook.
//
// Завдання 2: реалізуйте функцію SortByYear, яка сортує книги
// за PublishedYear без сторонніх бібліотек (sort.Slice або ручний
// алгоритм — обидва варіанти приймаються).
package library

// Author представляє автора книги.
//
// TODO (Завдання 1): додайте/скоригуйте поля на свій розсуд —
// мають бути щонайменше Name і Born.
type Author struct {
	Name string
	Born int
}

// Book представляє одну книгу в бібліотеці.
//
// Book.Author — приклад вкладеної структури, як і вимагає завдання.
//
// TODO (Завдання 1): додайте/скоригуйте поля на свій розсуд —
// мають бути щонайменше Title, Author і PublishedYear.
type Book struct {
	Title         string
	Author        Author
	PublishedYear int
}

// Library зберігає колекцію книг.
//
// Library.Books — поле-зріз, як і вимагає завдання.
type Library struct {
	Name  string
	Books []Book
}

// AddBook додає книгу до бібліотеки.
//
// TODO (Завдання 1): реалізуйте додавання b до l.Books.
func (l *Library) AddBook(b Book) {
	// TODO: ваш код тут
}

// SortByYear сортує books за PublishedYear (за зростанням) на місці,
// без використання сторонніх бібліотек.
//
// TODO (Завдання 2): реалізуйте сортування.
// Підказка: sort.Slice(books, func(i, j int) bool { ... }) —
// це вже частина стандартної бібліотеки Go, тому дозволена.
func SortByYear(books []Book) {
	// TODO: ваш код тут
}
