package into_struct

import "reflect"

// validateInto does some basic checks to help users of this class avoid common pitfalls with more helpful messages
func validateInto(dst interface{}) (rootV reflect.Value, err error) {
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
