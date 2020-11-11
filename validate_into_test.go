package into_struct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidateDestination(t *testing.T) {
	fixture := tWithSettableNestedStruct{}
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
			_, err := validateInto(c.input)
			if c.expected != nil {
				assert.EqualError(t, err, c.expected.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnmarshallInvalidDestination(t *testing.T) {
	assert.Error(t, Unmarshall(nil, new(mockParser)))
}
