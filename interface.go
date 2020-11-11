package into_struct

import "reflect"

type ValueSetter interface {
	// SetValue asks the callee to update the struct field in fieldV
	// fieldV is always Settable and it is the actual reflect.Value of the field, not a pointer to the field (unless the field is a pointer...)
	// You can also call fieldV.Addr().Interface() to get a pointer to the field, too
	// If you set the value, set handled = true. This will prevent the Unmarshall process from trying the default handler
	// return a non-nil err if you couldn't parse string value to the fieldV type. This will halt Unmarshalling.
	SetValue(structFullPath string, fieldV reflect.Value, structField reflect.StructField) (handled bool, err error)
}

type SliceLengther interface {
	// SliceLen returns how many elements are in the slice at the structFullPath
	// Zero is perfectly acceptable
	// err is any error encountered when looking up how many elements are in the array. This error will halt processing
	// and bubble up through to the caller of Unmarshall
	SliceLen(structFullPath string, fieldV reflect.Value, structField reflect.StructField) (length int, err error)
}

type Parser interface {
	ValueSetter
	SliceLengther
}
