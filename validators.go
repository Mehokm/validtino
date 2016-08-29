package validtino

// type val struct {
// 	NotEmpty ValidatorFunc
// 	Min      ValidatorFunc
// 	Between  ValidatorFunc
// 	Email    ValidatorFunc
// }
//
// func newVal() val {
// 	return val{
// 		NotEmpty: func(candidate interface{}, params interface{}) bool {
// 			switch candidate.(type) {
// 			case int:
// 				return candidate.(int) != 0
// 			case string:
// 				return candidate.(string) != ""
// 			default:
// 				return false
// 			}
// 		},
// 		Min: func(candidate interface{}, params interface{}) bool {
// 			min := params.Get("min").(int)
// 			switch candidate.(type) {
// 			case int:
// 				return candidate.(int) >= min
// 			case string:
// 				return utf8.RuneCountInString(candidate.(string)) >= min
// 			default:
// 				return false
// 			}
// 		},
// 		Between: func(candidate interface{}, params interface{}) bool {
// 			min := params.Get("min").(int)
// 			max := params.Get("max").(int)
// 			switch candidate.(type) {
// 			case int:
// 				return candidate.(int) >= min && candidate.(int) <= max
// 			case string:
// 				return utf8.RuneCountInString(candidate.(string)) >= min && utf8.RuneCountInString(candidate.(string)) <= max
// 			default:
// 				return false
// 			}
// 		},
// 		Email: func(candidate interface{}, params Painterface{}) bool {
// 			switch candidate.(type) {
// 			case string:
// 				match, _ := regexp.MatchString(`[-a-z0-9~!$%^&*_=+}{\'?]+(\.[-a-z0-9~!$%^&*_=+}{\'?]+)*@([a-z0-9_][-a-z0-9_]*(\.[-a-z0-9_]+)*\.(aero|arpa|biz|com|coop|edu|gov|info|int|mil|museum|name|net|org|pro|travel|mobi|[a-z][a-z])|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,5})?`, candidate.(string))
// 				return match
// 			default:
// 				return false
// 			}
// 		},
// 	}
// }
