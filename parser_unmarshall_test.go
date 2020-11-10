package into_struct

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type parserMock struct {
	mock.Mock
}

func (m *parserMock) SetValue(structFullPath string, fieldV reflect.Value) (handled bool, err error) {
	args := m.Called(structFullPath, fieldV)
	return args.Bool(0), args.Error(1)
}
func (m *parserMock) SliceLen(structFullPath string) (length int, err error) {
	args := m.Called(structFullPath)
	return args.Int(0), args.Error(1)
}

func TestUnmarshallStruct(t *testing.T) {
	cases := map[string]struct {
		mock        *parserMock
		expectedErr error
	}{
		"empty Parser": {
			mock: func() (m *parserMock) {
				m = &parserMock{}
				m.On("SetValue", "Settable", mock.Anything).
					Return(true, nil)
				return
			}(),
		},
		"value parse error": {
			mock: func() (m *parserMock) {
				m = &parserMock{}
				m.On("SetValue", "Settable", mock.Anything).
					Return(true, fmt.Errorf("parse error"))
				return
			}(),
			expectedErr: fmt.Errorf("parse error"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			var actual parserTestStruct
			err := Unmarshall(&actual, c.mock)
			if c.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, c.expectedErr.Error())
			}
			c.mock.AssertExpectations(t)
		})
	}
}

func TestUnmarshallInvalidDestination(t *testing.T) {
	err := Unmarshall(nil, new(parserMock))
	assert.Error(t, err, "")
}

type structWithSlice struct {
	Strings []string
}

func TestUnmarshallSlice(t *testing.T) {
	cases := map[string]struct {
		mock        *parserMock
		expectedErr error
	}{
		"empty slice": {
			mock: func() (m *parserMock) {
				m = &parserMock{}
				m.On("SliceLen", "Strings").
					Return(0, nil)
				return
			}(),
		},
		"single item slice": {
			mock: func() (m *parserMock) {
				m = &parserMock{}
				m.On("SliceLen", "Strings").
					Return(1, nil)
				m.On("SetValue", "Strings[0]", mock.Anything).
					Return(true, nil)
				return
			}(),
		},
		"slice parse error": {
			mock: func() (m *parserMock) {
				m = &parserMock{}
				m.On("SliceLen", "Strings").
					Return(-1, fmt.Errorf("parse error"))
				return
			}(),
			expectedErr: fmt.Errorf("parse error"),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			var actual structWithSlice
			err := Unmarshall(&actual, c.mock)
			if c.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, c.expectedErr.Error())
			}
			c.mock.AssertExpectations(t)
		})
	}
}

type structInnerStruct struct {
	Value string
}

type structWithStruct struct {
	Thing structInnerStruct
}

func TestUnmarshallStructInStruct(t *testing.T) {
	cases := map[string]struct {
		mock        *parserMock
		expectedErr error
	}{
		"parse inner struct": {
			mock: func() (m *parserMock) {
				m = &parserMock{}
				m.On("SetValue", "Thing", mock.Anything).
					Return(false, nil)
				m.On("SetValue", "Thing.Value", mock.Anything).
					Return(true, nil)
				return
			}(),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			var actual structWithStruct
			err := Unmarshall(&actual, c.mock)
			if c.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, c.expectedErr.Error())
			}
			c.mock.AssertExpectations(t)
		})
	}
}
