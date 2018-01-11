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
	validatorMap   map[string]*Validator
	structMap      map[string][]*property
	allowedTypeMap map[reflect.Kind]bool
	rParam, _      = regexp.Compile(`\([A-Za-z0-9,=' ]+\)`)
	mutex          = &sync.Mutex{}
)

// ValidatorFunc is the type of func that will be used to validate in your validator
type ValidatorFunc func(candidate, paramType interface{}) bool

// Validator is the type that is required to register a validator.  Name is the name of the validator - it matches the string in the tag
// Func is the function that is called to do the validation
// ParamType is required for mapping your validator parameters to your validator func
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

func init() {
	validatorMap = make(map[string]*Validator)
	structMap = make(map[string][]*property)
	allowedTypeMap = map[reflect.Kind]bool{
		reflect.Bool:          false,
		reflect.Int:           true,
		reflect.Int8:          true,
		reflect.Int16:         true,
		reflect.Int32:         true,
		reflect.Int64:         true,
		reflect.Uint:          true,
		reflect.Uint8:         true,
		reflect.Uint16:        true,
		reflect.Uint32:        true,
		reflect.Uint64:        true,
		reflect.Float32:       true,
		reflect.Float64:       true,
		reflect.Complex64:     false,
		reflect.Complex128:    false,
		reflect.Array:         false,
		reflect.Chan:          false,
		reflect.Func:          false,
		reflect.Interface:     true,
		reflect.Map:           false,
		reflect.Ptr:           false,
		reflect.Slice:         false,
		reflect.String:        true,
		reflect.Struct:        false,
		reflect.UnsafePointer: false,
	}

	RegisterValidator(NewContainsValidator())
	RegisterValidator(NewNotEmptyValidator())
	RegisterValidator(NewMinValidator())
	RegisterValidator(NewNumRangeValidator())
	RegisterValidator(NewEmailValidator())
}

// RegisterValidator allows a user to register a validator to use with validtino
func RegisterValidator(val *Validator) {
	mutex.Lock()
	defer mutex.Unlock()

	validatorMap[val.Name] = val
}

// RegisterStruct will speed up reflection for struct validation since it happened at start time
func RegisterStruct(s interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	sv := reflect.ValueOf(s)

	if sv.Kind() != reflect.Ptr {
		return errors.New("validtino: candidate must be ptr")
	}

	if sv.Elem().Kind() != reflect.Struct {
		return errors.New("validtino: candidate must be of type struct")
	}

	structMap[getKey(sv)] = getProperties(sv)

	return nil
}

// Validate will validate struct fields which have the valid tag and with
// corresponding validator
func Validate(s interface{}) []error {
	mutex.Lock()
	defer mutex.Unlock()

	var errs []error

	sv := reflect.ValueOf(s)

	if sv.Kind() != reflect.Ptr {
		return append(errs, errors.New("validtino: candidate must be ptr"))
	}

	sve := sv.Elem()

	if sve.Kind() != reflect.Struct {
		return append(errs, errors.New("validtino: candidate must be of type struct"))
	}

	key := getKey(sv)

	var props []*property
	var ok bool

	if props, ok = structMap[key]; !ok {
		props = getProperties(sv)
	}

	for _, prop := range props {
		prop.value = sve.FieldByName(prop.name).Interface()

		if len(prop.validatorParams) > 0 {
			setParamType(prop)
		}

		for _, vName := range prop.validatorNames {
			val := validatorMap[vName]
			passed := val.Func(prop.value, val.ParamType)

			if !passed {
				// check validator for custom message.  This could be the default (not implemented yet)
				err := fmt.Errorf("validtino: field '%s' failed validator '%s' with value '%v'", prop.name, val.Name, prop.value)
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func setParamType(prop *property) {
	for k, vName := range prop.validatorNames {
		val := validatorMap[vName]

		ptv := reflect.ValueOf(val.ParamType).Elem()

		numFields := ptv.NumField()
		for i := 0; i < numFields; i++ {
			ptField := ptv.Field(i)
			param := prop.validatorParams[k][i]

			switch ptField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var p int64
				pp, err := strconv.Atoi(param)
				if err == nil {
					p = int64(pp)
				}
				ptField.SetInt(p)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				var p uint64
				pp, err := strconv.Atoi(param)
				if err == nil {
					p = uint64(pp)
				}
				ptField.SetUint(p)
			case reflect.Float32:
				var p float64
				pp, err := strconv.ParseFloat(param, 32)
				if err == nil {
					p = float64(pp)
				}
				ptField.SetFloat(p)
			case reflect.Float64:
				var p float64
				pp, err := strconv.ParseFloat(param, 64)
				if err == nil {
					p = pp
				}
				ptField.SetFloat(p)
			case reflect.String, reflect.Interface:
				// check to see if value is single quoted for syntax sake
				// if it is not, set it to empty string
				// if it is, then remove the single quote

				if param[0] != 39 || param[len(param)-1] != 39 {
					param = ""
				} else {
					param = param[1 : len(param)-1]
				}
				ptField.SetString(param)
			}
		}
	}
}

func getProperties(sv reflect.Value) []*property {
	var props []*property

	sve := sv.Elem()
	numFields := sve.NumField()
	for i := 0; i < numFields; i++ {
		if !allowedTypeMap[sve.Field(i).Kind()] {
			continue
		}

		tField := sve.Type().Field(i)

		tag := tField.Tag.Get("valid")

		if tag == "" {
			continue
		}

		valTags := strings.Split(strings.Replace(tag, " ", "", -1), ";")

		var vNames []string
		var vParams [][]string

		for _, v := range valTags {
			vName := getValidatorName(v)

			if _, ok := validatorMap[vName]; !ok {
				continue
			}

			vNames = append(vNames, vName)

			rawParams := rParam.FindString(v)

			if rawParams != "" {
				tParams := strings.Split(rawParams[1:len(rawParams)-1], ",")
				vParams = append(vParams, tParams)
			}
		}

		prop := &property{
			name:            tField.Name,
			validatorNames:  vNames,
			validatorParams: vParams,
		}

		props = append(props, prop)
	}

	return props
}

func getKey(sv reflect.Value) string {
	sve := sv.Elem()

	return sve.Type().PkgPath() + "." + sve.Type().Name()
}

func getValidatorName(s string) string {
	i := strings.Index(s, "(")

	if i >= 0 {
		return s[:i]
	}

	return s
}
