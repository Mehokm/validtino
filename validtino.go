package validtino

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sync"
)

var (
	validatorMap map[string]ValidatorFunc
	rName, _     = regexp.Compile(`[A-Za-z]+`)
	rParam, _    = regexp.Compile(`\([A-Za-z0-9,=' ]+\)`)
	mutex        sync.RWMutex
)

type Params struct {
	params map[string]interface{}
}

func (p Params) Get(index string) interface{} {
	if v, ok := p.params[index]; ok {
		return v
	}

	return nil
}

type ValidatorFunc func(candidate interface{}, params Params) bool

func init() {
	validatorMap = make(map[string]ValidatorFunc)

	RegisterValidator(newVal())
}

func RegisterValidator(v interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return errors.New("validator must be of type struct")
	}

	numFields := rv.NumField()
	for i := 0; i < numFields; i++ {
		name := rv.Type().Field(i).Name
		funcVal := rv.Field(i).Interface().(ValidatorFunc)

		validatorMap[name] = funcVal
	}

	return nil
}

func RegisterValidatorFunc(name string, f ValidatorFunc) {
	mutex.Lock()
	defer mutex.Unlock()

	validatorMap[name] = f
}

// Validate will validate struct fields which have the valid tag and with
// corressponding validator
func Validate(s interface{}) []error {
	var errs []error

	sv := reflect.ValueOf(s)

	if sv.Kind() != reflect.Ptr {
		return append(errs, errors.New("candidate must be ptr"))
	}

	if sv.Elem().Kind() != reflect.Struct {
		return append(errs, errors.New("candidate must be of type struct"))
	}

	fields := getStructFields(sv)

	fmt.Println(fields)

	return errs
}

func getStructFields(sv reflect.Value) []reflect.Value {
	var fields []reflect.Value

	numFields := sv.Elem().NumField()
	for i := 0; i < numFields; i++ {
		field := sv.Elem().Field(i)
		tag := sv.Elem().Type().Field(i).Tag.Get("valid")

		if tag == "" {
			continue
		}

		validatorName := tag

		fmt.Println(field)
		fmt.Println(validatorName)
	}

	return fields
}

// type Validation struct {
// 	validators []interface{}
// }
//
// var rName, _ = regexp.Compile(`[A-Za-z]+`)
// var rParam, _ = regexp.Compile(`\([A-Za-z0-9,=' ]+\)`)
//
// func NewValidation(validators []interface{}) Validation {
// 	v := Validation{validators}
// 	return v
// }
//
// func (val Validation) Validate(s interface{}) (bool, []error) {
// 	passed := true
// 	sType := reflect.TypeOf(s)
// 	sValue := reflect.ValueOf(s)
//
// 	if sType.Kind() != reflect.Struct {
// 		return false, []error{errors.New("candidate must be of type Struct")}
// 	}
//
// 	numField := sType.NumField()
// 	errArr := []error{}
//
// 	for i := 0; i < numField; i++ {
// 		tField := sType.Field(i)
// 		vField := sValue.Field(i)
//
// 		valName := rName.FindString(tField.Tag.Get("valid"))
//
// 		if valName == "" {
// 			continue
// 		}
//
// 		var validator reflect.Value
// 		for _, v := range val.validators {
// 			validator = reflect.ValueOf(v).MethodByName(valName)
// 			if validator.IsValid() {
// 				break
// 			}
// 		}
//
// 		if !validator.IsValid() {
// 			return false, []error{errors.New(fmt.Sprintf("validator: '%v' is not defined", valName))}
// 		}
//
// 		valParamList := rParam.FindString(tField.Tag.Get("valid"))
//
// 		var valParams []string
// 		if valParamList != "" {
// 			valParamList = stripWhitespaceNotQuoted(valParamList[1 : len(valParamList)-1])
// 			valParams = strings.Split(valParamList, ",")
// 		}
//
// 		keyWords := make(map[string]string)
// 	L: // Restart the loop once you mutatue valParams
// 		for i, v := range valParams {
// 			if strings.Count(v, "=") == 1 {
// 				pair := strings.Split(v, "=")
// 				keyWords[pair[0]] = pair[1]
// 				valParams = append(valParams[:i], valParams[i+1:]...)
//
// 				goto L
// 			}
// 		}
//
// 		includeVal, includeOk := keyWords["include"]
// 		if includeOk {
// 			var includeValTester string
// 			switch vField.Kind() {
// 			case reflect.Int:
// 				includeValTester = strconv.Itoa(int(vField.Int()))
// 			case reflect.String:
// 				includeValTester = vField.String()
// 			}
// 			if includeVal == includeValTester {
// 				continue
// 			}
// 		}
//
// 		shouldExclude := false
// 		excludeVal, excludeOk := keyWords["exclude"]
// 		if excludeOk {
// 			var excludeValTester string
// 			switch vField.Kind() {
// 			case reflect.Int:
// 				excludeValTester = strconv.Itoa(int(vField.Int()))
// 			case reflect.String:
// 				excludeValTester = vField.String()
// 			}
// 			if excludeVal == excludeValTester {
// 				shouldExclude = true
// 			}
// 		}
//
// 		numParam := validator.Type().NumIn()
//
// 		if len(valParams) != numParam-1 {
// 			return false, []error{errors.New(fmt.Sprintf("validator: '%v' has a parameter mismatch. Providing %v, require %v",
// 				valName, len(valParams), numParam-1))}
// 		}
//
// 		inputs := make([]reflect.Value, numParam)
// 		// first param is always candidate
// 		inputs[0] = vField
//
// 		for i := 1; i < numParam; i++ {
// 			switch validator.Type().In(i).Kind() {
// 			case reflect.Int:
// 				p, _ := strconv.Atoi(valParams[i-1])
// 				inputs[i] = reflect.ValueOf(p)
// 			case reflect.String:
// 				inputs[i] = reflect.ValueOf(valParams[i-1])
// 			}
// 		}
//
// 		result := validator.Call(inputs)
// 		if shouldExclude || !result[0].Bool() {
// 			passed = false
//
// 			switch vField.Kind() {
// 			case reflect.Int:
// 				errArr = append(errArr,
// 					errors.New(fmt.Sprintf("field: '%v' failed validation with value %v. %v", tField.Name, vField.Int(), keyWords["message"])))
// 			case reflect.String:
// 				errArr = append(errArr,
// 					errors.New(fmt.Sprintf("field: '%v' failed validation with value '%v'. %v", tField.Name, vField.String(), keyWords["message"])))
// 			}
// 		}
// 	}
//
// 	return passed, errArr
// }
//
// func stripWhitespaceNotQuoted(s string) string {
// 	var strippedS string
//
// 	isQuote := false
// 	for _, char := range s {
// 		if string(char) == "'" {
// 			isQuote = !isQuote
// 		}
// 		if !isQuote && string(char) == " " {
// 			strippedS += ""
// 		} else {
// 			strippedS += string(char)
// 		}
// 	}
// 	return strippedS
// }
