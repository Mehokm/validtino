package main

import (
	"fmt"
	"time"
	"validtino"
)

type Test struct {
	A string `valid:"Contains('he')"`
	B string
	C int    `valid:"Min(3); Range(4, 9)"`
	D uint   `valid:"Min(7)"`
	E string `valid:"NotEmpty"`
}

func main() {
	t := Test{"hello", "bye", 2, uint(8), ""}

	validtino.RegisterValidator(validtino.NewMinValidator())
	validtino.RegisterValidator(validtino.NewRangeValidator())
	validtino.RegisterValidator(validtino.NewNotEmptyValidator())
	validtino.RegisterValidator(validtino.NewContainsValidator())

	validtino.RegisterStruct(&Test{})

	start := time.Now()
	errs := validtino.Validate(&t)
	end := time.Now()

	fmt.Println(end.Sub(start))
	fmt.Println(errs)
}
