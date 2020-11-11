package into_struct

import (
	"reflect"
)

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func Unmarshall(into interface{}, parser Parser) (err error) {
	rootV, err := validateInto(into)
	if err != nil {
		return
	}
	rootVElem := rootV.Elem()
	err = unmarshallStruct(parser, "", rootVElem)
	return
}

// unmarshallStruct is the internal method, which can be called recursively. This performs the heavy-lifting
func unmarshallStruct(parser Parser, structParentPath string, structRefV reflect.Value) (err error) {
	for i := 0; i < structRefV.NumField(); i++ {
		fieldV := structRefV.Field(i)
		fieldT := structRefV.Type().Field(i)
		err = unmarshallField(parser, structParentPath, fieldV, fieldT)
		if err != nil {
			return
		}
	}
	return
}

// unmarshallField unmarshalls a value into a single field in a struct. Could be the root struct or a nested struct
func unmarshallField(parser Parser, structParentPath string, fieldV reflect.Value, fieldT reflect.StructField) (err error) {
	if fieldV.CanSet() {
		structFullPath := appendStructPath(&structParentPath, &fieldT.Name)
		switch fieldT.Type.Kind() {
		case reflect.Slice:
			err = unmarshallSlice(parser, structFullPath, fieldV, fieldT)
		default:
			err = unmarshallValue(parser, structFullPath, fieldV, fieldT)
		}
	}
	return
}

// unmarshallValue extracts a single value and sets it to a value in a struct
func unmarshallValue(parser Parser, structFullPath string, fieldV reflect.Value, fieldT reflect.StructField) (err error) {
	var wasSet bool
	wasSet, err = parser.SetValue(structFullPath, fieldV, fieldT)
	if err != nil {
		return
	}
	if wasSet {
		// Value set, no fall-back needed
		return
	}
	switch fieldT.Type.Kind() {
	case reflect.Struct:
		// fall back: no value found or was not set due to lack of type support
		err = unmarshallStruct(parser, structFullPath, fieldV)
	default:
		err = NewErrProgramming("unsupported fallback type for field: " + fieldT.PkgPath + "." + fieldT.Name + " only Struct and Slice are supported")
	}
	return
}

// unmarshallSlice operates on a slice of objects. It will initialize the slice, then populate all of its members
// from the environment variables
func unmarshallSlice(parser Parser, sliceFieldPath string, sliceValue reflect.Value, fieldT reflect.StructField) (err error) {
	var length int
	length, err = parser.SliceLen(sliceFieldPath, sliceValue, fieldT)
	if err != nil {
		return
	}
	if length > 0 {
		sliceValueType := sliceValue.Type()
		newSlice := reflect.MakeSlice(sliceValueType, length, length)
		sliceValue.Set(newSlice)
		for i := 0; i < length; i++ {
			sliceElement := newSlice.Index(i)
			sliceElementT := sliceElement.Type()
			sliceIndexPath := appendStructIndex(&sliceFieldPath, i)
			switch sliceElementT.Kind() {
			case reflect.Struct:
				err = unmarshallStruct(parser, sliceIndexPath, sliceElement)
			default:
				err = unmarshallValue(parser, sliceIndexPath, sliceElement, fieldT)
			}
			if err != nil {
				return
			}
		}
	}
	return
}
