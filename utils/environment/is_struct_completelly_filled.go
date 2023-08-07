package environment

import "reflect"

func IsStructCompletellyFilled(e any) (isCompletelyFilled bool, missingField string) {
	var (
		reflectValue reflect.Value
		fieldValue   reflect.Value
		fieldName    string
	)

	// Iterate over each field in the struct using reflection
	reflectValue = reflect.ValueOf(e).Elem()
	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue = reflectValue.Field(i)
		fieldName = reflectValue.Type().Field(i).Name

		// Check for []string
		if fieldValue.Kind() == reflect.Slice && fieldValue.Type().Elem().Kind() == reflect.String {
			if fieldValue.Len() == 0 {
				isCompletelyFilled = false
				missingField = fieldName
				return
			}
			continue
		}

		// Check if the field is uninitialized (has zero value)
		if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
			isCompletelyFilled = false
			missingField = fieldName
			return // Return false if any field is undefined
		}
	}
	// Return true if all fields are defined
	isCompletelyFilled = true

	return
}

type IsStructCompletellyFilledF func(e any) (isCompletelyFilled bool, missingField string)

func IsStructNotCompletelyFilled(e any) (isNotCompletelyFilled bool, missingField string) {
	var (
		isCompletelyFilled bool
	)

	isCompletelyFilled, missingField = IsStructCompletellyFilled(e)
	isNotCompletelyFilled = !isCompletelyFilled

	return

}

type IsStructNotCompletelyFilledF func(e any) (isNotCompletelyFilled bool, missingField string)
