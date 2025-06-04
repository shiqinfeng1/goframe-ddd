package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUidIsValid(t *testing.T) {
	tests := []struct {
		name     string
		uid      string
		expected bool
	}{
		{
			name:     "valid uid with MAC",
			uid:      "prefix-00:1A:2B:3C:4D:5E",
			expected: true,
		},
		{
			name:     "invalid uid - no delimiter",
			uid:      "noDelimiter",
			expected: false,
		},
		{
			name:     "invalid uid - empty string",
			uid:      "",
			expected: false,
		},
		{
			name:     "invalid uid - only prefix",
			uid:      "prefix-",
			expected: false,
		},
		{
			name:     "invalid uid - invalid MAC format",
			uid:      "prefix-invalidMAC",
			expected: false,
		},
		{
			name:     "invalid uid - incomplete MAC",
			uid:      "prefix-00:1A:2B",
			expected: false,
		},
		{
			name:     "valid uid with lowercase MAC",
			uid:      "prefix-00:1a:2b:3c:4d:5e",
			expected: true,
		},
		{
			name:     "invalid uid with hyphenated MAC",
			uid:      "prefix-00-1A-2B-3C-4D-5E",
			expected: false,
		},
		{
			name:     "invalid uid - multiple delimiters but invalid MAC",
			uid:      "prefix1-prefix2-invalidMAC",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UidIsValid(tt.uid)
			assert.Equal(t, tt.expected, result)
		})
	}
}
