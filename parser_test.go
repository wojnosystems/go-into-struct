package into_struct

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type tStructWithField struct {
	Value string
}

type tWithSettableNestedStruct struct {
	Settable tStructWithField `name:"settable"`
}

type tWithSlice struct {
	Strings []string
}

type tWithSliceWithStruct struct {
	Strings []tStructWithField
}

type tWithUnsupportedType struct {
	Chan chan int
}

func TestUnmarshall(t *testing.T) {
	cases := map[string]struct {
		into        interface{}
		mock        *mockParser
		expectedErr error
	}{
		"field parsed": {
			into: &tStructWithField{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SetValue", "Value", mock.Anything, mock.Anything).
					Return(true, nil)
				return
			}(),
		},
		"field is unsupported": {
			into: &tWithUnsupportedType{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SetValue", "Chan", mock.Anything, mock.Anything).
					Return(false, nil)
				return
			}(),
			expectedErr: NewErrProgramming("unsupported fallback type for field: .Chan only Struct and Slice are supported"),
		},
		"nested parse Settable as a value": {
			into: &tWithSettableNestedStruct{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SetValue", "Settable", mock.Anything, mock.Anything).
					Return(true, nil)
				return
			}(),
		},
		"nested parse Settable as a struct": {
			into: &tWithSettableNestedStruct{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SetValue", "Settable", mock.Anything, mock.Anything).
					Return(false, nil)
				m.On("SetValue", "Settable.Value", mock.Anything, mock.Anything).
					Return(true, nil)
				return
			}(),
		},
		"nested value parse error": {
			into: &tWithSettableNestedStruct{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SetValue", "Settable", mock.Anything, mock.Anything).
					Return(true, fmt.Errorf("parse error"))
				return
			}(),
			expectedErr: fmt.Errorf("parse error"),
		},
		"slice empty": {
			into: &tWithSlice{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SliceLen", "Strings", mock.Anything, mock.Anything).
					Return(0, nil)
				return
			}(),
		},
		"slice single item": {
			into: &tWithSlice{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SliceLen", "Strings", mock.Anything, mock.Anything).
					Return(1, nil)
				m.On("SetValue", "Strings[0]", mock.Anything, mock.Anything).
					Return(true, nil)
				return
			}(),
		},
		"slice multiple item": {
			into: &tWithSlice{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SliceLen", "Strings", mock.Anything, mock.Anything).
					Return(3, nil)
				m.On("SetValue", "Strings[0]", mock.Anything, mock.Anything).
					Return(true, nil)
				m.On("SetValue", "Strings[1]", mock.Anything, mock.Anything).
					Return(true, nil)
				m.On("SetValue", "Strings[2]", mock.Anything, mock.Anything).
					Return(true, nil)
				return
			}(),
		},
		"slice index parse error": {
			into: &tWithSlice{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SliceLen", "Strings", mock.Anything, mock.Anything).
					Return(-1, fmt.Errorf("parse error"))
				return
			}(),
			expectedErr: fmt.Errorf("parse error"),
		},
		"slice item parse error": {
			into: &tWithSlice{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SliceLen", "Strings", mock.Anything, mock.Anything).
					Return(1, nil)
				m.On("SetValue", "Strings[0]", mock.Anything, mock.Anything).
					Return(false, fmt.Errorf("parse error"))
				return
			}(),
			expectedErr: fmt.Errorf("parse error"),
		},
		"slice item parsed": {
			into: &tWithSlice{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SliceLen", "Strings", mock.Anything, mock.Anything).
					Return(1, nil)
				m.On("SetValue", "Strings[0]", mock.Anything, mock.Anything).
					Return(true, nil)
				return
			}(),
		},
		"slice with struct item parsed": {
			into: &tWithSliceWithStruct{},
			mock: func() (m *mockParser) {
				m = &mockParser{}
				m.On("SliceLen", "Strings", mock.Anything, mock.Anything).
					Return(1, nil)
				m.On("SetValue", "Strings[0].Value", mock.Anything, mock.Anything).
					Return(true, nil)
				return
			}(),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			err := Unmarshall(c.into, c.mock)
			if c.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, c.expectedErr.Error())
			}
			c.mock.AssertExpectations(t)
		})
	}
}
