package into_struct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type parserInnerTestStruct struct {
	Value string
}

type parserTestStruct struct {
	Settable   parserInnerTestStruct `name:"settable"`
	unsettable parserInnerTestStruct
}

func Test_ValidateDestination(t *testing.T) {
	fixture := parserTestStruct{}
	cases := map[string]struct {
		input    interface{}
		expected error
	}{
		"working": {
			input: &fixture.Settable,
		},
		"not a struct": {
			input:    &fixture.Settable.Value,
			expected: NewErrProgramming(`'into' argument must be a struct`),
		},
		"not settable: pass by value": {
			input:    fixture.Settable,
			expected: NewErrProgramming(`'into' argument must be a reference`),
		},
		"nil": {
			input:    nil,
			expected: NewErrProgramming(`'into' argument must be not be nil`),
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			_, err := validateDestination(c.input)
			if c.expected != nil {
				assert.EqualError(t, err, c.expected.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_AppendStructPath(t *testing.T) {
	cases := map[string]struct {
		input    string
		parent   string
		expected string
	}{
		"empty": {},
		"no parent": {
			input:    "test",
			expected: "test",
		},
		"with parent": {
			input:    "test",
			parent:   "parent",
			expected: "parent.test",
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := appendStructPath(c.parent, c.input)
			assert.Equal(t, c.expected, actual)
		})
	}
}

//func Test_FieldNameOrDefault(t *testing.T) {
//	cases := map[string]struct{
//
//	}{
//
//	}
//
//	for caseName, c := range cases {
//		t.Run(caseName, func(t *testing.T) {
//
//		})
//	}
//}
