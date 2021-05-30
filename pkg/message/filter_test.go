package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFilterSingleErrors(t *testing.T) {
	cases := []struct {
		name   string
		filter string
		errMsg string
	}{
		{
			name:   "Invalid characters",
			filter: "module=m/",
			errMsg: ErrorFilterInvalidChars,
		},
		{
			name:   "Inversion in middle",
			filter: "mmodule!=m",
			errMsg: ErrorFilterInvertNotAtBeginning,
		},
		{
			name:   "Multiple fields",
			filter: "module=m=m",
			errMsg: ErrorFilterInvalidFields,
		},
		{
			name:   "No fields",
			filter: "module",
			errMsg: ErrorFilterInvalidFields,
		},
		{
			name:   "Invalid key",
			filter: "m=m",
			errMsg: ErrorFilterInvalidKey,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkFilterSingle(tc.filter)
			assert.EqualError(t, err, tc.errMsg)
		})
	}
}

func TestCheckFilterErrors(t *testing.T) {
	cases := []struct {
		name   string
		filter string
		resNil bool
	}{
		{
			name:   "Multiple valid filters",
			filter: "module=m+!dest=d",
			resNil: true,
		},
		{
			name:   "One valid filter",
			filter: "module=m",
			resNil: true,
		},
		{
			name:   "One invalid filter",
			filter: "module!=m",
			resNil: false,
		},
		{
			name:   "One invalid filter, one correct",
			filter: "module=m/+dest=d",
			resNil: false,
		},
		{
			name:   "Multiple invalid filters",
			filter: "module=m/+dest!=d",
			resNil: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CheckFilter(tc.filter)
			assert.Equal(t, tc.resNil, got == nil)
		})
	}
}

func TestFilterSingleReturnsError(t *testing.T) {
	_, err := filterSingle(&Message{}, "module!=m")
	assert.NotNil(t, err)
}

var filters = []struct {
	name    string
	message Message
	filter  string
	result  bool
}{
	{
		name: "Filter module",
		message: Message{
			Module: "mod",
		},
		filter: "module=mod",
		result: true,
	},
	{
		name: "Filter module with inversion",
		message: Message{
			Module: "mod",
		},
		filter: "!module=mod",
		result: false,
	},
	{
		name: "Filter command",
		message: Message{
			Command: "test",
		},
		filter: "command=test",
		result: true,
	},
	{
		name: "False condition",
		message: Message{
			Module: "mod",
		},
		filter: "module=othermod",
		result: false,
	},
	{
		name: "False condition with inversion",
		message: Message{
			Module: "mod",
		},
		filter: "!module=mod",
		result: false,
	},
}

func TestCheckFilterSingleValid(t *testing.T) {
	for _, tc := range filters {
		t.Run(tc.name, func(t *testing.T) {
			err := checkFilterSingle(tc.filter)
			assert.Nil(t, err)
		})
	}
}

func TestFilterMessageSingle(t *testing.T) {
	for _, tc := range filters {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := filterSingle(&tc.message, tc.filter)

			assert.Equal(t, tc.result, got)
		})
	}
}

func TestInvert(t *testing.T) {
	cases := []struct {
		name     string
		b        bool
		inv      bool
		expected bool
	}{
		{
			name:     "Invert true",
			b:        true,
			inv:      true,
			expected: false,
		},
		{
			name:     "Don't invert true",
			b:        true,
			inv:      false,
			expected: true,
		},
		{
			name:     "Invert false",
			b:        false,
			inv:      true,
			expected: true,
		},
		{
			name:     "Don't invert false",
			b:        false,
			inv:      false,
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := invert(tc.b, tc.inv)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestFilter(t *testing.T) {
	cases := []struct {
		name    string
		message Message
		filter  string
		result  bool
	}{
		{
			name: "One true",
			message: Message{
				Module: "mod",
			},
			filter: "module=mod",
			result: true,
		},
		{
			name: "Multiple true",
			message: Message{
				Module:  "mod",
				Command: "test",
			},
			filter: "module=mod+command=test",
			result: true,
		},
		{
			name: "One false",
			message: Message{
				Module: "mod",
			},
			filter: "module=othermod",
			result: false,
		},
		{
			name: "One false, one true",
			message: Message{
				Module:  "mod",
				Command: "test",
			},
			filter: "module=othermod+command=test",
			result: false,
		},
		{
			name: "One inverted, one normal",
			message: Message{
				Module:  "mod",
				Command: "test",
			},
			filter: "!module=othermod+command=test",
			result: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Filter(&tc.message, tc.filter)

			assert.Nil(t, err)
			assert.Equal(t, tc.result, got)
		})
	}
}
