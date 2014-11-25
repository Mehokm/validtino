package main

import (
	"fmt"
	"validatino"
)

type TestStruct struct {
	StringNotEmpty string `valid:"NotEmpty"`
	StringEmpty    string `valid:"NotEmpty"`
	Size           int    `valid:"Min(2, message = 'This will be a custom message that will display')"`
	Range          int    `valid:"Range(2, 10, exclude = 3)"`
	Email          string `valid:"Email(include = root)"`
	Tester         string `valid:"MyTester"`
}

type TestValidator struct{}

func (tv TestValidator) MyTester(candidate interface{}) bool {
	return false
}

func main() {
	tsFail := TestStruct{"Hello", "", 1, 3, "root", "This should fail because MyTester validator returns false"}

	var myVal validatino.Validator
	var myTestVal TestValidator

	val := validatino.NewValidation([]interface{}{myVal, myTestVal})

	passed, errors := val.Validate(tsFail)
	if !passed && len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}
