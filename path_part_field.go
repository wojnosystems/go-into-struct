package into_struct

type pathPartField struct {
	pathField
}

func (s *pathPartField) Name() string {
	return s.name
}

func (s *pathPartField) String() string {
	return s.Name()
}
