package into_struct

import (
	"reflect"
	"strconv"
)

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func Unmarshall(into interface{}, parser Parser) (err error) {
	rootV, err := validateDestination(into)
	if err != nil {
		return
	}
	rootVElem := rootV.Elem()
	rootTElem := rootVElem.Type()
	err = unmarshallStruct(parser, "", rootVElem, rootTElem)
	return
}

// validateDestination does some basic checks to help users of this class avoid common pitfalls with more helpful messages
func validateDestination(dst interface{}) (rootV reflect.Value, err error) {
	if dst == nil {
		err = NewErrProgramming("'into' argument must be not be nil")
		return
	}
	rootV = reflect.ValueOf(dst)
	rootT := rootV.Type()
	if rootT.Kind() != reflect.Ptr {
		err = NewErrProgramming("'into' argument must be a reference")
		return
	}
	if rootV.Elem().Kind() != reflect.Struct {
		err = NewErrProgramming("'into' argument must be a struct")
		return
	}
	return
}

// unmarshallStruct is the internal method, which can be called recursively. This performs the heavy-lifting
func unmarshallStruct(parser Parser, structParentPath string, structRefV reflect.Value, structRefT reflect.Type) (err error) {
	for i := 0; i < structRefV.NumField(); i++ {
		fieldV := structRefV.Field(i)
		fieldT := structRefT.Field(i)
		err = unmarshallField(parser, structParentPath, fieldV, fieldT, structRefT)
		if err != nil {
			return
		}
	}
	return
}

// unmarshallField unmarshalls a value into a single field in a struct. Could be the root struct or a nested struct
func unmarshallField(parser Parser, structParentPath string, fieldV reflect.Value, fieldT reflect.StructField, parentT reflect.Type) (err error) {
	if fieldV.CanSet() {
		structFullPath := appendStructPath(structParentPath, fieldT.Name)

		if fieldT.Type.Kind() == reflect.Slice {
			err = unmarshallSlice(parser, structFullPath, fieldV)
		} else {
			err = unmarshallValue(parser, structFullPath, fieldV, fieldT.Type)
		}
		if err != nil {
			return
		}
	}
	return
}

// unmarshallValue extracts a single value and sets it to a value in a struct
func unmarshallValue(parser Parser, structFullPath string, fieldV reflect.Value, fieldT reflect.Type) (err error) {
	var wasSet bool
	wasSet, err = parser.SetValue(structFullPath, fieldV)
	if err != nil {
		return
	}
	if wasSet {
		// Value set, no fall-back needed
		return
	}
	if fieldT.Kind() == reflect.Struct {
		// fall back: no value found or was not set due to lack of type support
		err = unmarshallStruct(parser, structFullPath, fieldV, fieldT)
	}
	return
}

// appendStructPath concatenates the parent path name with the current field's name
func appendStructPath(parent string, name string) string {
	if parent != "" {
		return parent + "." + name
	}
	return name
}

// unmarshallSlice operates on a slice of objects. It will initialize the slice, then populate all of its members
// from the environment variables
func unmarshallSlice(parser Parser, sliceFieldPath string, sliceValue reflect.Value) (err error) {
	var length int
	length, err = parser.SliceLen(sliceFieldPath)
	if err != nil {
		return
	}
	if length > 0 {
		sliceValueType := sliceValue.Type()
		newSlice := reflect.MakeSlice(sliceValueType, length, length)
		sliceValue.Set(newSlice)
		for i := 0; i < length; i++ {
			sliceElement := newSlice.Index(i)
			sliceElementType := sliceElement.Type()
			indexPath := sliceFieldPath + "[" + strconv.FormatInt(int64(i), 10) + "]"
			err = unmarshallValue(parser, indexPath, sliceElement, sliceElementType)
			if err != nil {
				return
			}
		}
	}
	return
}
