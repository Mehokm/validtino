package main

import (
	"fmt"
	"time"

	"github.com/Mehokm/validtino"
)

type Test struct {
	A string `valid:"Contains('he')"`
	B string
	C int    `valid:"Min(3); NumRange(4, 9)"`
	D uint   `valid:"Min(7)"`
	E string `valid:"NotEmpty"`
	F string
	G string `valid:"Contains('usi')"`
	H string `valid:"NotEmpty"`
	I string `valid:"Contains('wh')"`
	J string `valid:"Contains('wo')"`
	K string `valid:"NotEmpty"`
}

func main() {
	t := Test{"hello", "bye", 2, uint(8), "", "", "using", "s", "what", "work", ""}

	validtino.RegisterStruct(&Test{})

	start := time.Now()
	errs := validtino.Validate(&t)
	end := time.Now()

	fmt.Println(end.Sub(start))
	fmt.Println(errs)
}
