package main

import (
	"fmt"
	"validtino"
)

type Test struct {
	A string `valid:"NotEmpty"`
	B string
}

func main() {
	t := Test{"I am not empty", "Also not empty, but don't care"}

	errs := validtino.Validate(&t)

	if len(errs) > 0 {
		fmt.Println(errs)
	}
}
