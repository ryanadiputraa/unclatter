package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	cases := []struct {
		name     string
		arg      string
		expected bool
	}{
		{
			name:     "should return true for a valid email address",
			arg:      "test@mail.com",
			expected: true,
		},
		{
			name:     "should return false for a invalid email address",
			arg:      "testmail",
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			isValid := IsValidEmail(c.arg)
			assert.Equal(t, c.expected, isValid)
		})
	}
}
