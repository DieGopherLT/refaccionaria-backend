package validator

import (
	"reflect"
)

// HasEmptyStringField checks if a string field on a given struct is empty
func HasEmptyStringField(object interface{}) bool {
	fields := reflect.TypeOf(object)
	values := reflect.ValueOf(object)
	num := fields.NumField()

	for i := 0; i < num; i++ {
		value := values.Field(i)
		if value.Kind() == reflect.String {
			v := value.String()
			if v == "" {
				return true
			}
		} else {
			continue
		}
	}
	return false
}
