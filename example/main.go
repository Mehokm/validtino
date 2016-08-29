package main

import (
	"fmt"
	"reflect"
	"unicode/utf8"
	"validtino"
)

type Test struct {
	A string `valid:"NotEmpty"`
	B string
	C int `valid:"Min(3)"`
}

type MinParamType struct {
	Min int
}

func main() {
	v := validtino.Validator{
		Name: "min",
		Func: func(candidate interface{}, t interface{}) bool {
			param := t.(*MinParamType)
			switch candidate.(type) {
			case int:
				return candidate.(int) >= param.Min
			case string:
				return utf8.RuneCountInString(candidate.(string)) >= param.Min
			default:
				return false
			}
		},
		ParamType: new(MinParamType),
	}

	pt := reflect.ValueOf(v.ParamType).Elem()

	pt.Field(0).SetInt(10)

	b := v.Func(11, v.ParamType)

	fmt.Println(b)
}
