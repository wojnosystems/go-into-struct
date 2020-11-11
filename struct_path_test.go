package into_struct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			actual := appendStructPath(&c.parent, &c.input)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func Test_AppendStructIndex(t *testing.T) {
	cases := map[string]struct {
		input    int
		parent   string
		expected string
	}{
		"no parent": {
			input:    6,
			expected: "[6]",
		},
		"with parent": {
			input:    0,
			parent:   "parent",
			expected: "parent[0]",
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := appendStructIndex(&c.parent, c.input)
			assert.Equal(t, c.expected, actual)
		})
	}
}
