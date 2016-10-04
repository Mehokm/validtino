package main

import (
	"fmt"
	"unicode/utf8"
	"validtino"
)

type Test struct {
	A string `valid:"Min(5)"`
	B string
	C int `valid:"Min(3); Range(4, 9)"`
}

type NumParamType struct {
	Num int
}

type RangeParamType struct {
	Low, High int
}

func main() {
	v := &validtino.Validator{
		Name:      "Min",
		ParamType: NumParamType{},
		Func: func(candidate interface{}, t interface{}) bool {
			param := t.(NumParamType)
			switch candidate.(type) {
			case int:
				return candidate.(int) >= param.Num
			case string:
				return utf8.RuneCountInString(candidate.(string)) >= param.Num
			default:
				return false
			}
		},
	}

	v2 := &validtino.Validator{
		Name:      "Range",
		ParamType: RangeParamType{},
		Func: func(candidate interface{}, t interface{}) bool {
			param := t.(RangeParamType)
			switch candidate.(type) {
			case int:
				return candidate.(int) >= param.Low && candidate.(int) <= param.High
			case string:
				return utf8.RuneCountInString(candidate.(string)) >= param.Low &&
					utf8.RuneCountInString(candidate.(string)) <= param.High
			default:
				return false
			}
		},
	}

	t := Test{"hello", "bye", 2}

	validtino.RegisterValidator(v)
	validtino.RegisterValidator(v2)

	errs := validtino.Validate(&t)

	fmt.Println(errs)
}
