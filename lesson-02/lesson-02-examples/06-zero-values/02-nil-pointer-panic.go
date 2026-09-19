// Розділ 4.4 · Пастки nil
//
// Демонструє: розіменування nil-вказівника завжди панікує.
// Використано recover(), щоб демонстрація не завершувала термінал
// аварійно — панічне повідомлення просто друкується.
package main

import "fmt"

type User struct {
	Name string
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("сталася паніка:", r)
		}
	}()

	var u *User
	fmt.Println(u.Name) // паніка: розіменування nil-вказівника
}
