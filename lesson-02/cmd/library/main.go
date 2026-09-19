// Команда library будує приклад бібліотеки та виводить її через
// fmt.Printf("%+v", ...), як і вимагає Завдання 1.
package main

import (
	"fmt"

	"example.com/lesson02/library"
)

func main() {
	lib := library.Library{Name: "City Library"}

	lib.AddBook(library.Book{
		Title:         "The Go Programming Language",
		Author:        library.Author{Name: "Alan Donovan", Born: 1971},
		PublishedYear: 2015,
	})
	lib.AddBook(library.Book{
		Title:         "The C Programming Language",
		Author:        library.Author{Name: "Dennis Ritchie", Born: 1941},
		PublishedYear: 1978,
	})

	library.SortByYear(lib.Books)

	fmt.Printf("%+v\n", lib)
}
