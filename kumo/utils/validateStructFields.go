package utils

import "reflect"

func IsStructFilledCompletely(s any) (isFilledCompletely bool) {
	var (
		reflectValue reflect.Value
		fieldValue   reflect.Value
	)

	// Iterate over each field in the struct using reflection
	reflectValue = reflect.ValueOf(s).Elem()
	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue = reflectValue.Field(i)

		// Check if the field is uninitialized (has zero value)
		if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
			return false // Return false if any field is undefined
		}
	}
	// Return true if all fields are defined
	return true
}
