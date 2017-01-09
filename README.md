Validtino
==========

Validtino was created in order to provide a simple way to validate structs in go.  Inspired by the validator provided by Hibernate, Validtino uses struct field tags to provide an easy way to define any struct field with an appropriate validator.

In addition to the built in validators, it is easy to create and use your own!  Following the example below, all you need to do is to create a \*Validator type with a Name, ParamType (the type you you use in your validator function), and Func (of ValidatorFunc type)

The ValidatorFunc type will require two parameters, the candidate and the param type, and return bool.  The candidate and param type will be both of type interface{}.

The ParamType is required because it allows you to define a type that will be mapped to the parameters of the validator defined in the struct tag.  So as an example, the tag ````valid:NumRange(4, 9)"```` has the validator NumRange with parameters 4 and 9.

Now we will create a param type: 
```go
type NumRangeParamType struct {
	Low, High int
}
```
This param type will be used like so with out validator:
```go
&Validator{
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
```
So to finish up the example, the number 4 will be mapped to Low and 9 will be mapped to High.  You can then use those properties on your param type in your validation function

The param type will be specific to your validator, but you can use one to more than one validator if it satisfies the validator parameters.  You can name it whatever you want, becasue you will be type converting it from interface to your param type as in the above example.

To register your validators with Validtino, you will use the RegisterValidator function.

When calling Validate, you will get a slice of errors as the result.  If everything passed, the slice will be of len == 0; 

Example Usage:
==============

```go
package main

import (
	"fmt"
	"strings"
	"validtino"
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
	K string `valid:"CustomVal('foo')"`
}

type CustomValParamType struct {
	String string
}

func main() {
	validtino.RegisterValidator(customVal())

	t := Test{"hello", "bye", 2, uint(8), "", "", "using", "s", "what", "work", "bar"}

	errs := validtino.Validate(&t)

	for _, err := range errs {
		fmt.Println(err)
	}
}

func customVal() *validtino.Validator {
	return &validtino.Validator{
		Name:      "CustomVal",
		ParamType: CustomValParamType{},
		Func: func(candidate interface{}, t interface{}) bool {
			s := t.(CustomValParamType)
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
```

This will produce the following output:

```
validtino: field 'C' failed validator 'Min' with value '2'
validtino: field 'C' failed validator 'NumRange' with value '2'
validtino: field 'E' failed validator 'NotEmpty' with value ''
validtino: field 'K' failed validator 'CustomVal' with value 'bar'
```
