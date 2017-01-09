package validtino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Foo int `valid:"Test(4, 19)"`
}

type TestStructMissing struct {
	Foo int `valid:""`
}

type TestStructInvalidType struct {
	Foo []int `valid:"Test(4, 19)"`
}

type TestStructNoValidator struct {
	Foo int `valid:"Test2(1, 10)"`
}

type TestParamType struct {
	A int
	B int
}

func TestRegisterValidator(t *testing.T) {
	// given
	val := getTestVal()

	// when
	RegisterValidator(val)

	// then
	assert.Equal(t, validatorMap["Test"], val)
}

func TestRegisterStruct(t *testing.T) {
	// given
	s := TestStruct{}

	// when
	RegisterStruct(&s)

	// then
	assert.NotEmpty(t, structMap["validtino.TestStruct"])
}

func TestRegisterStruct_mustBePtr(t *testing.T) {
	// given
	s := TestStruct{}

	// when
	err := RegisterStruct(s)

	// then
	assert.EqualError(t, err, "validtino: candidate must be ptr")
}

func TestRegisterStruct_mustBeStruct(t *testing.T) {
	// given
	var s interface{}

	// when
	err := RegisterStruct(&s)

	// then
	assert.EqualError(t, err, "validtino: candidate must be of type struct")
}

func TestValidate_validatesStructCorrectlyNoError(t *testing.T) {
	// given
	RegisterValidator(getTestVal())

	s := TestStruct{6}

	// when
	errs := Validate(&s)

	// then
	assert.Empty(t, errs)
}

func TestValidate_validatesStructCorrectlyWithError(t *testing.T) {
	// given
	RegisterValidator(getTestVal())

	s := TestStruct{1}

	// when
	errs := Validate(&s)

	// then
	assert.Equal(t, len(errs), 1)
	assert.EqualError(t, errs[0], "validtino: field 'Foo' failed validator 'Test' with value '1'")
}

func TestValidate_errorMustBePtr(t *testing.T) {
	// given
	RegisterValidator(getTestVal())

	s := TestStruct{1}

	// when
	errs := Validate(s)

	// then
	assert.Equal(t, len(errs), 1)
	assert.EqualError(t, errs[0], "validtino: candidate must be ptr")
}

func TestValidate_errorMustBeStruct(t *testing.T) {
	// given
	RegisterValidator(getTestVal())

	var s interface{}

	// when
	errs := Validate(&s)

	// then
	assert.Equal(t, len(errs), 1)
	assert.EqualError(t, errs[0], "validtino: candidate must be of type struct")
}

func TestValidate_ignoresMissingValidator(t *testing.T) {
	// given
	RegisterValidator(getTestVal())

	s := TestStructMissing{1}

	// when
	errs := Validate(&s)

	// then
	assert.Equal(t, len(errs), 0)
}

func TestValidate_ignoresInvalidType(t *testing.T) {
	// given
	RegisterValidator(getTestVal())

	s := TestStructInvalidType{[]int{1}}

	// when
	errs := Validate(&s)

	// then
	assert.Equal(t, len(errs), 0)
}

func TestValidate_ignoresUnregisteredValidator(t *testing.T) {
	// given
	RegisterValidator(getTestVal())

	s := TestStructNoValidator{1}

	// when
	errs := Validate(&s)

	// then
	assert.Equal(t, len(errs), 0)
}

func getTestVal() *Validator {
	return &Validator{
		Name:      "Test",
		ParamType: TestParamType{},
		Func: func(candidate interface{}, paramType interface{}) bool {
			p := paramType.(TestParamType)
			return (candidate.(int) >= p.A && candidate.(int) <= p.B)
		},
	}
}
