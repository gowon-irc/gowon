package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitOptArray(t *testing.T) {
	cases := []struct {
		name     string
		arr      []string
		expected []string
	}{
		{
			name:     "no changes",
			arr:      []string{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "split one",
			arr:      []string{"hello,world", "goodbye"},
			expected: []string{"hello", "world", "goodbye"},
		},
		{
			name:     "split all",
			arr:      []string{"hello,world", "goodbye,world"},
			expected: []string{"hello", "world", "goodbye", "world"},
		},
		{
			name:     "empty list",
			arr:      []string{},
			expected: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := splitOptArray(tc.arr)

			assert.Equal(t, got, tc.expected)
		})
	}
}
