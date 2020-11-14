package into_struct

import "fmt"

type pathPartSlice struct {
	pathField
	index int
}

func (s *pathPartSlice) Name() string {
	return s.name
}

func (s *pathPartSlice) String() string {
	return fmt.Sprintf(`%s[%d]`, s.name, s.index)
}

func (s *pathPartSlice) Index() int {
	return s.index
}
