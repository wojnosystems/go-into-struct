package into_struct

import "reflect"

type ValueSetter interface {
	// SetValue asks the callee to update the struct field in fieldV
	// fieldV is always Settable and it is the actual reflect.Value of the field, not a pointer to the field (unless the field is a pointer...)
	// You can also call fieldV.Addr().Interface() to get a pointer to the field, too
	// If you set the value, set handled = true. This will prevent the Unmarshall process from trying the default handler
	// return a non-nil err if you couldn't parse string value to the fieldV type. This will halt Unmarshalling.
	SetValue(structFullPath string, fieldV reflect.Value) (handled bool, err error)
}

type SliceLengther interface {
	// SliceLen returns how many elements are in the slice at the structFullPath
	// Zero is perfectly acceptable
	// err is any error encountered when looking up how many elements are in the array. This error will halt processing
	// and bubble up through to the caller of Unmarshall
	SliceLen(structFullPath string) (length int, err error)
}

type FieldNamer interface {
	// FieldName provides the parser the ability to specify the canonical name that a field in the destination struct may have. Tags are usually used for this purpose. For example, the field may be defined as: "Name string `tagName:"special"`" In this case, tagName's value of "special" could be used as the structFullPath name for this field.
	// structParentPath is provided for convenience and represents the path to the struct within which the fieldT exists.
	// structParent is the type of struct within which the fieldT exists.
	// fieldT is the reflect.StructField representing the field being Unmarshalled.
	// return a non-empty string with the official name of the field. If you wish to use the default name, which
	// matches what exists in code for your struct, simply return an empty string.
	FieldName(structParentPath string, structParent reflect.Type, fieldT reflect.StructField) (fieldName string)
}

type Parser interface {
	ValueSetter
	SliceLengther
	FieldNamer
}
