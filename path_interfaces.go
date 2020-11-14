package into_struct

import "reflect"

type PathParter interface {
	Value() reflect.Value
	Type() reflect.Type
	StructField() reflect.StructField
	Name() string
	String() string
}

type PathSliceParter interface {
	PathParter
	Index() int
}

type pathField struct {
	name         string
	reflectValue reflect.Value
	reflectType  reflect.Type
	reflectField reflect.StructField
}

func pathFieldFromParter(parter PathParter) pathField {
	return pathField{
		name:         parter.StructField().Name,
		reflectValue: parter.Value(),
		reflectType:  parter.Type(),
		reflectField: parter.StructField(),
	}
}

func (b pathField) Type() reflect.Type {
	return b.reflectType
}

func (b pathField) Value() reflect.Value {
	return b.reflectValue
}

func (b pathField) StructField() reflect.StructField {
	return b.reflectField
}
