Validatino
==========

Validatino was created in order to provide a simple way to validate structs in go.  Inspired by the validator provided by Hibernate, Validatino uses struct field tags to provide an easy way to define any struct field with an appropriate validator.

In addition to the built in validators, it is easy to create and use your own!  Following the example below, all you need to do is create your validator type and use it as a recevier to your validator methods.  

Validator methods can contain any number of parameters that are required, but the first parameter will always contain the candidate (value of the struct field) that will be tested.  Order will matter (except in the case of keywords) of the parameters, so the method signature should match the signature when you apply the validator.  For example, the Range validator will take a min and max value.  So, the defined method will have parameters: candidate, min, max.

To register your validators with Validatino, pass them in as an array -- simple as that.

Calling the Validate method on your struct var will run each valid validator defined on the fields.  It will return whether the struct passed or failed, and if there were any errors.  

As of now, there are only three supported keywords: include, exclude, and message.  Include will "include" whatever value set to it and pass if the field value equals that include value.  Exclude is the opposite, if the field value equals the exclude value then it will return false for that field validation.  Message is just a custom error message that you can defined for when field validation fails

Example Usage:
==============

```go
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
```

This will produce the following output:

```
field: 'StringEmpty' failed validation with value ''.
field: 'Size' failed validation with value 1. 'This will be a custom message that will display'
field: 'Range' failed validation with value 3.
field: 'Tester' failed validation with value 'This should fail because MyTester validator returns false'.
```
