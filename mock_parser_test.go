package into_struct

import (
	"github.com/stretchr/testify/mock"
	"reflect"
)

type mockParser struct {
	mock.Mock
}

func (m *mockParser) SetValue(structFullPath string, fieldV reflect.Value, fieldT reflect.StructField) (handled bool, err error) {
	args := m.Called(structFullPath, fieldV, fieldT)
	return args.Bool(0), args.Error(1)
}

func (m *mockParser) SliceLen(structFullPath string, fieldV reflect.Value, structField reflect.StructField) (length int, err error) {
	args := m.Called(structFullPath, fieldV, structField)
	return args.Int(0), args.Error(1)
}
