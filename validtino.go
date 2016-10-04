package validtino

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	validatorMap map[string]*Validator
	rName, _     = regexp.Compile(`[A-Za-z]+`)
	rParam, _    = regexp.Compile(`\([A-Za-z0-9,=' ]+\)`)
	mutex        sync.RWMutex
)

type Validator struct {
	Name      string
	Func      ValidatorFunc
	ParamType interface{}
}

type property struct {
	name            string
	value           interface{}
	validatorNames  []string
	validatorParams [][]string
}

type ValidatorFunc func(candidate interface{}, paramType interface{}) bool

func init() {
	validatorMap = make(map[string]*Validator)
}

// RegisterValidator allows a user to register a validator to use with validtino
func RegisterValidator(val *Validator) {
	mutex.Lock()
	defer mutex.Unlock()

	validatorMap[val.Name] = val
}

// Validate will validate struct fields which have the valid tag and with
// corresponding validator
func Validate(s interface{}) []error {
	var errs []error

	sv := reflect.ValueOf(s)

	if sv.Kind() != reflect.Ptr {
		return append(errs, errors.New("candidate must be ptr"))
	}

	if sv.Elem().Kind() != reflect.Struct {
		return append(errs, errors.New("candidate must be of type struct"))
	}

	props := getProperties(sv)

	for _, prop := range props {
		// add check to see if this is necessary.  They might want to define this
		// way upstream
		setParamType(prop)

		for _, vName := range prop.validatorNames {
			val := validatorMap[vName]

			passed := val.Func(prop.value, val.ParamType)
			if !passed {
				err := fmt.Errorf("field: '%v' failed validator '%v' with value %v", prop.name, val.Name, prop.value)
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func setParamType(prop property) {
	for k, vName := range prop.validatorNames {
		val := validatorMap[vName]

		ptCopy := reflect.New(reflect.TypeOf(val.ParamType)).Elem()

		numFields := ptCopy.NumField()
		for i := 0; i < numFields; i++ {
			ptField := ptCopy.Field(i)
			param := prop.validatorParams[k][i]

			switch ptField.Kind() {
			case reflect.Int:
				var p int64
				pp, err := strconv.Atoi(param)
				if err == nil {
					p = int64(pp)
				}
				ptField.SetInt(p)
			case reflect.String:
				ptField.SetString(param)
			}
		}

		val.ParamType = ptCopy.Interface()
	}
}

func getProperties(sv reflect.Value) []property {
	mutex.RLock()
	defer mutex.RUnlock()

	var props []property

	se := sv.Elem()
	numFields := se.NumField()
	for i := 0; i < numFields; i++ {
		field := se.Field(i)
		tField := se.Type().Field(i)

		tag := tField.Tag.Get("valid")

		if tag == "" {
			continue
		}

		valTags := strings.Split(strings.Replace(tag, " ", "", -1), ";")

		var vNames []string
		var vParams [][]string

		for _, v := range valTags {
			vName := rName.FindString(v)

			if _, ok := validatorMap[vName]; !ok {
				continue
			}

			vNames = append(vNames, vName)
			vParams = append(vParams, strings.Split(strings.Trim(rParam.FindString(v), "()"), ","))
		}

		prop := property{
			name:            tField.Name,
			value:           field.Interface(),
			validatorNames:  vNames,
			validatorParams: vParams,
		}

		props = append(props, prop)
	}

	return props
}
