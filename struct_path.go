package into_struct

import "fmt"

// appendStructPath concatenates the parent path name with the current field's name
// parent can be an empty string, but cannot be null. parent is not modified. Pointer used to avoid a copy
// name is the field, struct, or slice field name to append to the parent name. name is not modified. Pointer used to avoid a copy
// returns the field name appended to the parent path as a new string
func appendStructPath(parent *string, fieldName *string) string {
	if *parent != "" {
		return fmt.Sprintf("%s.%s", *parent, *fieldName)
	}
	return *fieldName
}

// appendStructIndex concatenates the parent path name with the current field's slice index
// parent can be an empty string, but cannot be null. parent is not modified. Pointer used to avoid a copy
// index is the integer offset index of the item in the slice
// returns the field name appended to the parent path as a new string
func appendStructIndex(parent *string, index int) string {
	return fmt.Sprintf("%s[%d]", *parent, index)
}
