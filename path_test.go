package into_struct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPath_Top(t *testing.T) {
	cases := map[string]struct {
		input    Path
		expected PathParter
	}{
		"empty": {
			input: func() Path {
				return Path{parts: []PathParter{}}
			}(),
		},
		"only top": {
			input: func() Path {
				return Path{parts: []PathParter{
					&pathPartField{pathField{name: "top"}},
				}}
			}(),
			expected: &pathPartField{pathField{name: "top"}},
		},
		"top with parent": {
			input: func() Path {
				return Path{parts: []PathParter{
					&pathPartField{pathField{name: "parent"}},
					&pathPartField{pathField{name: "top"}},
				}}
			}(),
			expected: &pathPartField{pathField{name: "top"}},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.input.Top()
			if c.expected == nil {
				assert.Nil(t, actual)
			} else {
				assert.Equal(t, c.expected.String(), actual.String())
			}
		})
	}
}

func TestPath_ParentOfTop(t *testing.T) {
	cases := map[string]struct {
		input    Path
		expected PathParter
	}{
		"empty": {
			input: func() Path {
				return Path{parts: []PathParter{}}
			}(),
		},
		"only top": {
			input: func() Path {
				return Path{parts: []PathParter{
					&pathPartField{pathField{name: "top"}},
				}}
			}(),
		},
		"top with parent": {
			input: func() Path {
				return Path{parts: []PathParter{
					&pathPartField{pathField{name: "parent"}},
					&pathPartField{pathField{name: "top"}},
				}}
			}(),
			expected: &pathPartField{pathField{name: "parent"}},
		},
		"top with parent and others": {
			input: func() Path {
				return Path{parts: []PathParter{
					&pathPartField{pathField{name: "parent3"}},
					&pathPartField{pathField{name: "parent2"}},
					&pathPartField{pathField{name: "parent"}},
					&pathPartField{pathField{name: "top"}},
				}}
			}(),
			expected: &pathPartField{pathField{name: "parent"}},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := c.input.ParentOfTop()
			if c.expected == nil {
				assert.Nil(t, actual)
			} else {
				assert.Equal(t, c.expected.String(), actual.String())
			}
		})
	}
}
