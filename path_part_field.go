package into_struct

type pathPartField struct {
	pathField
}

func (s pathPartField) String() string {
	return s.name
}
