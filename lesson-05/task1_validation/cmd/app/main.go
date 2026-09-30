// Command app demonstrates registration validation and reporting all invalid fields.
package main

import (
	"errors"
	"fmt"

	"lesson05/task1_validation"
)

func main() {
	form := validation.RegistrationForm{
		Email:    "",
		Password: "short",
		Age:      -1,
	}

	err := validation.ValidateRegistration(form)
	if err == nil {
		fmt.Println("registration is valid")
		return
	}

	var validationErr *validation.ValidationError
	if errors.As(err, &validationErr) {
		fmt.Println("invalid registration fields:", validationErr.Fields)
		return
	}

	fmt.Println("unexpected validation error:", err)
}
