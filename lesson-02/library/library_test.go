package library

import "testing"

// TestAddBook перевіряє Завдання 1: додавання книги до бібліотеки
// та коректність вкладеної структури Book.Author.
func TestAddBook(t *testing.T) {
	lib := Library{Name: "City Library"}

	b := Book{
		Title:         "The Go Programming Language",
		Author:        Author{Name: "Alan Donovan", Born: 1971},
		PublishedYear: 2015,
	}

	lib.AddBook(b)

	if len(lib.Books) != 1 {
		t.Fatalf("len(lib.Books) = %d, want 1", len(lib.Books))
	}
	if lib.Books[0].Title != b.Title {
		t.Errorf("lib.Books[0].Title = %q, want %q", lib.Books[0].Title, b.Title)
	}
	if lib.Books[0].Author.Name != "Alan Donovan" {
		t.Errorf("lib.Books[0].Author.Name = %q, want %q", lib.Books[0].Author.Name, "Alan Donovan")
	}
}

func TestAddBook_Multiple(t *testing.T) {
	lib := Library{Name: "City Library"}

	lib.AddBook(Book{Title: "Book A", PublishedYear: 2001})
	lib.AddBook(Book{Title: "Book B", PublishedYear: 2002})
	lib.AddBook(Book{Title: "Book C", PublishedYear: 2003})

	if len(lib.Books) != 3 {
		t.Fatalf("len(lib.Books) = %d, want 3", len(lib.Books))
	}
}

// TestSortByYear — таблично-керований тест для Завдання 2.
// Перевіряє вже відсортований, зворотно відсортований і порожній зріз.
func TestSortByYear(t *testing.T) {
	tests := []struct {
		name  string
		input []Book
		want  []int // очікувана послідовність PublishedYear після сортування
	}{
		{
			name: "уже відсортований",
			input: []Book{
				{Title: "A", PublishedYear: 1999},
				{Title: "B", PublishedYear: 2005},
				{Title: "C", PublishedYear: 2020},
			},
			want: []int{1999, 2005, 2020},
		},
		{
			name: "відсортований у зворотному порядку",
			input: []Book{
				{Title: "A", PublishedYear: 2020},
				{Title: "B", PublishedYear: 2005},
				{Title: "C", PublishedYear: 1999},
			},
			want: []int{1999, 2005, 2020},
		},
		{
			name:  "порожній зріз",
			input: []Book{},
			want:  []int{},
		},
		{
			name: "один елемент",
			input: []Book{
				{Title: "Solo", PublishedYear: 2010},
			},
			want: []int{2010},
		},
		{
			name: "однакові роки",
			input: []Book{
				{Title: "A", PublishedYear: 2000},
				{Title: "B", PublishedYear: 1990},
				{Title: "C", PublishedYear: 2000},
			},
			want: []int{1990, 2000, 2000},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			SortByYear(tt.input)

			if len(tt.input) != len(tt.want) {
				t.Fatalf("len after sort = %d, want %d", len(tt.input), len(tt.want))
			}
			for i, b := range tt.input {
				if b.PublishedYear != tt.want[i] {
					t.Errorf("position %d: PublishedYear = %d, want %d", i, b.PublishedYear, tt.want[i])
				}
			}
		})
	}
}

// TestSortByYear_SortsInPlace перевіряє, що сортування відбувається
// на місці (in place) — SortByYear не повинна повертати нову
// послідовність окремо від вхідного зрізу.
func TestSortByYear_SortsInPlace(t *testing.T) {
	books := []Book{
		{Title: "Z", PublishedYear: 2020},
		{Title: "A", PublishedYear: 2000},
	}
	original := books // той самий базовий масив

	SortByYear(books)

	if original[0].Title != "A" {
		t.Errorf("сортування має бути in-place: original[0].Title = %q, want %q", original[0].Title, "A")
	}
}
