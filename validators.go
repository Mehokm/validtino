package validtino

import (
	"strings"
	"unicode/utf8"
)

type NumParamType struct {
	Num int
}

type StringParamType struct {
	Str string
}

type RangeParamType struct {
	Low, High int
}

func NewContainsValidator() *Validator {
	return &Validator{
		Name:      "Contains",
		ParamType: StringParamType{},
		Func: func(candidate interface{}, t interface{}) bool {
			s := t.(StringParamType)
			switch candidate.(type) {
			case int:
				return false
			case string:
				return strings.Contains(candidate.(string), s.Str)
			default:
				return false
			}
		},
	}
}

func NewNotEmptyValidator() *Validator {
	return &Validator{
		Name: "NotEmpty",
		Func: func(candidate interface{}, t interface{}) bool {
			switch candidate.(type) {
			case int:
				return candidate.(int) > 0
			case string:
				return utf8.RuneCountInString(candidate.(string)) > 0
			default:
				return false
			}
		},
	}
}

func NewMinValidator() *Validator {
	return &Validator{
		Name:      "Min",
		ParamType: NumParamType{},
		Func: func(candidate interface{}, t interface{}) bool {
			param := t.(NumParamType)
			switch candidate.(type) {
			case int:
				return candidate.(int) >= param.Num
			case uint:
				return candidate.(uint) >= uint(param.Num)
			case string:
				return utf8.RuneCountInString(candidate.(string)) >= param.Num
			default:
				return false
			}
		},
	}
}

func NewRangeValidator() *Validator {
	return &Validator{
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
}
