package validtino

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type NumParamType struct {
	Number int
}

type StringParamType struct {
	String string
}

type NumRangeParamType struct {
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
				return strings.Contains(candidate.(string), s.String)
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
				return candidate.(int) >= param.Number
			case uint:
				return candidate.(uint) >= uint(param.Number)
			case string:
				return utf8.RuneCountInString(candidate.(string)) >= param.Number
			default:
				return false
			}
		},
	}
}

func NewNumRangeValidator() *Validator {
	return &Validator{
		Name:      "NumRange",
		ParamType: NumRangeParamType{},
		Func: func(candidate interface{}, t interface{}) bool {
			param := t.(NumRangeParamType)
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

func NewEmailValidator() *Validator {
	return &Validator{
		Name: "Email",
		Func: func(candidate interface{}, t interface{}) bool {
			switch candidate.(type) {
			case string:
				// regex from http://emailregex.com
				valid, _ := regexp.MatchString(
					"^(([^<>()\\[\\]\\.,;:\\s@\"]+(\\.[^<>()\\[\\]\\.,;:\\s@\"]+)*)|(\".+\"))@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}])|(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2,}))$",
					candidate.(string),
				)

				return valid
			default:
				return false
			}
		},
	}
}
