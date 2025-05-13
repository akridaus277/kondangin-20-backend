package utils

import (
	"reflect"
)

// IsNilOrEmpty returns true if the value is nil or empty (zero value for its type)
func IsNilOrEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}

	return false
}
