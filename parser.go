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
	err = unmarshallStruct(parser, Path{}, rootVElem)
	return
}

// unmarshallStruct is the internal method, which can be called recursively. This performs the heavy-lifting
func unmarshallStruct(parser Parser, structParentPath Path, structRefV reflect.Value) (err error) {
	for i := 0; i < structRefV.NumField(); i++ {
		fieldV := structRefV.Field(i)
		fieldT := structRefV.Type().Field(i)
		structParentPath.within(&pathPartField{
			pathField: pathField{
				name:         fieldT.Name,
				reflectValue: fieldV,
				reflectType:  fieldT.Type,
				reflectField: fieldT,
			},
		}, func() {
			err = unmarshallField(parser, structParentPath)
		})
		if err != nil {
			return
		}
	}
	return
}

// unmarshallField unmarshalls a value into a single field in a struct. Could be the root struct or a nested struct
func unmarshallField(parser Parser, structParentPath Path) (err error) {
	if structParentPath.Top().Value().CanSet() {
		switch structParentPath.Top().Type().Kind() {
		case reflect.Slice:
			err = unmarshallSlice(parser, structParentPath)
		default:
			err = unmarshallValue(parser, structParentPath)
		}
	}
	return
}

// unmarshallValue extracts a single value and sets it to a value in a struct
func unmarshallValue(parser Parser, structFullPath Path) (err error) {
	var wasSet bool
	wasSet, err = parser.SetValue(structFullPath)
	if err != nil {
		return
	}
	if wasSet {
		// Value set, no fall-back needed
		return
	}
	top := structFullPath.Top()
	switch top.Type().Kind() {
	case reflect.Struct:
		// fall back: no value found or was not set due to lack of type support
		err = unmarshallStruct(parser, structFullPath, structFullPath.Top().Value())
	default:
		err = NewErrProgramming("unsupported fallback type for field: " + top.Type().PkgPath() + "." + top.String() + " only Struct and Slice are supported")
	}
	return
}

// unmarshallSlice operates on a slice of objects. It will initialize the slice, then populate all of its members
// from the environment variables
func unmarshallSlice(parser Parser, sliceFieldPath Path) (err error) {
	var length int
	length, err = parser.SliceLen(sliceFieldPath)
	if err != nil {
		return
	}
	if length > 0 {
		sliceField := sliceFieldPath.Top()
		sliceIndexPath := &pathPartSlice{
			pathField: pathFieldFromParter(sliceField),
			index:     0,
		}
		sliceIndexPath.reflectType = sliceField.Type().Elem()
		sliceFieldPath.setTop(sliceIndexPath)
		newSlice := reflect.MakeSlice(sliceField.Type(), length, length)
		sliceField.Value().Set(newSlice)
		for i := 0; i < length; i++ {
			sliceElement := newSlice.Index(i)
			sliceIndexPath.reflectValue = sliceElement
			sliceIndexPath.index = i
			err = unmarshallField(parser, sliceFieldPath)
			if err != nil {
				return
			}
		}
	}
	return
}
