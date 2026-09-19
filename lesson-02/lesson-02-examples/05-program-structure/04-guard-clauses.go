// Розділ 3.4 · Найкращі практики
//
// Демонструє: ранній return (охоронні вирази) замість глибокої
// вкладеності — та сама логіка, дві версії для порівняння.
package main

import (
	"errors"
	"fmt"
)

type User struct {
	IsActive bool
	Balance  float64
}

// Гірше: вкладеність зростає з кожною перевіркою
func processNested(user *User) error {
	if user != nil {
		if user.IsActive {
			if user.Balance > 0 {
				fmt.Println("processing (nested version)...")
				return nil
			}
		}
	}
	return errors.New("cannot process")
}

// Краще: охоронні вирази, основна логіка на верхньому рівні
func processGuardClauses(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if !user.IsActive {
		return errors.New("user is not active")
	}
	if user.Balance <= 0 {
		return errors.New("insufficient balance")
	}
	fmt.Println("processing (guard-clause version)...")
	return nil
}

func main() {
	u := &User{IsActive: true, Balance: 100}

	if err := processNested(u); err != nil {
		fmt.Println("error:", err)
	}
	if err := processGuardClauses(u); err != nil {
		fmt.Println("error:", err)
	}

	bad := &User{IsActive: false, Balance: 100}
	if err := processGuardClauses(bad); err != nil {
		fmt.Println("expected error:", err)
	}
}
