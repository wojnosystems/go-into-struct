package into_struct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrProgramming_Error(t *testing.T) {
	cases := map[string]struct {
		msg      string
		expected string
	}{
		"random string": {
			msg:      "test",
			expected: "programming error: test",
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			assert.EqualError(t, NewErrProgramming(c.msg), c.expected)
		})
	}
}
