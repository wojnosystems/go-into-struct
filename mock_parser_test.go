package into_struct

import (
	"github.com/stretchr/testify/mock"
)

type mockParser struct {
	mock.Mock
}

func (m *mockParser) SetValue(structFullPath Path) (handled bool, err error) {
	args := m.Called(structFullPath.String())
	return args.Bool(0), args.Error(1)
}

func (m *mockParser) SliceLen(structFullPath Path) (length int, err error) {
	args := m.Called(structFullPath.String())
	return args.Int(0), args.Error(1)
}
