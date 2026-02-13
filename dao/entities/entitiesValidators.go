// Package entities provides type-safe wrappers and validation utilities for database types.
// This file contains reflection-based validators for struct fields and type checking.
package entities

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
)

// GetFunctionName returns the name of the function passed as an argument
// It uses reflection to obtain the function's program counter and retrieves its name.
func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

// GetStructType extracts and returns the type name of a struct from any data value.
// It handles pointers and fully-qualified type names (e.g., "package.Type" becomes "Type").
// Returns the type name wrapped as a Table.
func GetStructType(data any) Table {
	logHandler.Trace.Printf("Resolving Struct Type for data: %v", data)
	rtnType := reflect.TypeOf(data).String()
	base := rtnType

	// base := rtnType
	// If the type is a pointer, get the underlying type
	if strings.Contains(rtnType, "*") {
		rtnType = reflect.TypeOf(data).Elem().String()
	}
	// If the type is a struct, get the name of the struct
	if strings.Contains(rtnType, ".") {
		rtnType = strings.Split(rtnType, ".")[1]
	}
	logHandler.Trace.Printf("{TYPE} Resolved Struct Type: %v (base: %v)", rtnType, base)

	return Table(rtnType)
}

// IsValidFieldInStruct validates whether a field name exists in the given struct.
// It normalizes the input by unwrapping pointers and handling slices/arrays to get
// to the underlying struct type, then uses reflection to check for the field.
//
// Parameters:
//   - fromField: The field name to validate
//   - data: The struct instance (or pointer/slice thereof) to validate against
//
// Returns:
//   - nil if the field exists in the struct
//   - ErrInvalidFieldWrapper if the field doesn't exist or data is invalid
func IsValidFieldInStruct(fromField Field, data any) error {
	// Normalise the type: unwrap pointers, and if it's a slice/array, use the element type.
	if data == nil {
		logHandler.Error.Printf("Cannot validate field '%v' on <nil> data", fromField.String())
		return commonErrors.ErrInvalidFieldWrapper(fromField.String())
	}

	t := reflect.TypeOf(data)
	for t.Kind() == reflect.Ptr {
		// Unwrap pointer types
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		logHandler.Trace.Printf("Validating Field '%v' against slice/array element type '%v'", fromField.String(), t.Elem())
		t = t.Elem()
		for t.Kind() == reflect.Ptr {
			// Handle slices of pointers to structs
			t = t.Elem()
		}
	}

	logHandler.Trace.Printf("Validating Field '%v' in Struct type '%v'", fromField.String(), t.Name())

	if t.Kind() != reflect.Struct {
		logHandler.Error.Printf("Type '%v' is not a struct; cannot validate field '%v'", t, fromField.String())
		return commonErrors.ErrInvalidFieldWrapper(fromField.String())
	}

	if _, isValidField := t.FieldByName(fromField.String()); !isValidField {
		logHandler.Error.Printf("Field '%v' not found in struct '%v'", fromField.String(), t.Name())
		logHandler.Error.Println(commonErrors.ErrInvalidFieldWrapper(fromField.String()))
		return commonErrors.ErrInvalidFieldWrapper(fromField.String())
	}
	logHandler.Trace.Printf("Field '%v' is valid in struct '%v'", fromField.String(), t.Name())
	return nil
}

// IsValidTypeForField validates whether the provided data matches the expected type
// of a specific field in a struct. This ensures type safety when setting struct fields.
//
// Parameters:
//   - field: The field name to check
//   - data: The value whose type should match the field's type
//   - forStruct: The struct instance (or pointer/slice thereof) containing the field
//
// Returns:
//   - nil if the data type matches the field's expected type
//   - ErrInvalidFieldWrapper if the field doesn't exist
//   - ErrInvalidTypeWrapper if there's a type mismatch
//
// Example:
//
//	err := IsValidTypeForField(Field("Name"), "John", &User{})
//	// Returns nil if User.Name is of type string, error otherwise
func IsValidTypeForField(field Field, data, forStruct any) error {
	if forStruct == nil {
		logHandler.Error.Printf("Cannot validate type for field '%v' on <nil> struct", field.String())
		return commonErrors.ErrInvalidFieldWrapper(field.String())
	}

	// Normalise the type of forStruct: unwrap pointers, and if it's a slice/array, use the element type.
	st := reflect.TypeOf(forStruct)
	for st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	if st.Kind() == reflect.Slice || st.Kind() == reflect.Array {
		st = st.Elem()
		for st.Kind() == reflect.Ptr {
			st = st.Elem()
		}
	}

	if st.Kind() != reflect.Struct {
		logHandler.Error.Printf("Type '%v' is not a struct; cannot validate type for field '%v'", st, field.String())
		return commonErrors.ErrInvalidFieldWrapper(field.String())
	}

	structField, found := st.FieldByName(field.String())
	if !found {
		logHandler.Error.Printf("Field '%v' not found in struct '%v' when validating type", field.String(), st.Name())
		return commonErrors.ErrInvalidFieldWrapper(field.String())
	}

	dataType := "<nil>"
	if data != nil {
		dataType = reflect.TypeOf(data).String()
	}
	structFieldType := structField.Type.String()
	if dataType != structFieldType {
		logHandler.Error.Printf("Type mismatch for field '%v': expected '%v', got '%v'", field.String(), structFieldType, dataType)
		return commonErrors.ErrInvalidTypeWrapper(field.String(), dataType, structFieldType)
	}
	logHandler.Trace.Printf("Type for field '%v' is valid: expected '%v', got '%v'", field.String(), structFieldType, dataType)
	return nil
}
