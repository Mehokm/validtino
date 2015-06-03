package validatino

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Validation struct {
	validators []interface{}
}

var rName, _ = regexp.Compile(`[A-Za-z]+`)
var rParam, _ = regexp.Compile(`\([A-Za-z0-9,=' ]+\)`)

func NewValidation(validators []interface{}) Validation {
	v := Validation{validators}
	return v
}

func (val Validation) Validate(s interface{}) (bool, []error) {
	passed := true
	sType := reflect.TypeOf(s)
	sValue := reflect.ValueOf(s)

	if sType.Kind() != reflect.Struct {
		return false, []error{errors.New("candidate must be of type Struct")}
	}

	numField := sType.NumField()
	errArr := []error{}

	for i := 0; i < numField; i++ {
		tField := sType.Field(i)
		vField := sValue.Field(i)

		valName := rName.FindString(tField.Tag.Get("valid"))

		if valName == "" {
			continue
		}

		var validator reflect.Value
		for _, v := range val.validators {
			validator = reflect.ValueOf(v).MethodByName(valName)
			if validator.IsValid() {
				break
			}
		}

		if !validator.IsValid() {
			return false, []error{errors.New(fmt.Sprintf("validator: '%v' is not defined", valName))}
		}

		valParamList := rParam.FindString(tField.Tag.Get("valid"))

		var valParams []string
		if valParamList != "" {
			valParamList = stripWhitespaceNotQuoted(valParamList[1 : len(valParamList)-1])
			valParams = strings.Split(valParamList, ",")
		}

		keyWords := make(map[string]string)
	L: // Restart the loop once you mutatue valParams
		for i, v := range valParams {
			if strings.Count(v, "=") == 1 {
				pair := strings.Split(v, "=")
				keyWords[pair[0]] = pair[1]
				valParams = append(valParams[:i], valParams[i+1:]...)

				goto L
			}
		}

		includeVal, includeOk := keyWords["include"]
		if includeOk {
			var includeValTester string
			switch vField.Kind() {
			case reflect.Int:
				includeValTester = strconv.Itoa(int(vField.Int()))
			case reflect.String:
				includeValTester = vField.String()
			}
			if includeVal == includeValTester {
				continue
			}
		}

		shouldExclude := false
		excludeVal, excludeOk := keyWords["exclude"]
		if excludeOk {
			var excludeValTester string
			switch vField.Kind() {
			case reflect.Int:
				excludeValTester = strconv.Itoa(int(vField.Int()))
			case reflect.String:
				excludeValTester = vField.String()
			}
			if excludeVal == excludeValTester {
				shouldExclude = true
			}
		}

		numParam := validator.Type().NumIn()

		if len(valParams) != numParam-1 {
			return false, []error{errors.New(fmt.Sprintf("validator: '%v' has a parameter mismatch. Providing %v, require %v",
				valName, len(valParams), numParam-1))}
		}

		inputs := make([]reflect.Value, numParam)
		// first param is always candidate
		inputs[0] = vField

		for i := 1; i < numParam; i++ {
			switch validator.Type().In(i).Kind() {
			case reflect.Int:
				p, _ := strconv.Atoi(valParams[i-1])
				inputs[i] = reflect.ValueOf(p)
			case reflect.String:
				inputs[i] = reflect.ValueOf(valParams[i-1])
			}
		}

		result := validator.Call(inputs)
		if shouldExclude || !result[0].Bool() {
			passed = false

			switch vField.Kind() {
			case reflect.Int:
				errArr = append(errArr,
					errors.New(fmt.Sprintf("field: '%v' failed validation with value %v. %v", tField.Name, vField.Int(), keyWords["message"])))
			case reflect.String:
				errArr = append(errArr,
					errors.New(fmt.Sprintf("field: '%v' failed validation with value '%v'. %v", tField.Name, vField.String(), keyWords["message"])))
			}
		}
	}

	return passed, errArr
}

func stripWhitespaceNotQuoted(s string) string {
	var strippedS string

	isQuote := false
	for _, char := range s {
		if string(char) == "'" {
			isQuote = !isQuote
		}
		if !isQuote && string(char) == " " {
			strippedS += ""
		} else {
			strippedS += string(char)
		}
	}
	return strippedS
}
