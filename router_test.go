package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandRouterAdd(t *testing.T) {
	cr := &CommandRouter{}
	cmd := &Command{Command: "command"}
	cr.Add(cmd)

	assert.Equal(t, "command", cr.Commands[0].GetCommand())
}

func createCommandRouterFromPriorities(in []int) *CommandRouter {
	cr := CommandRouter{}

	for _, i := range in {
		cmd := &HttpCommand{
			Command:  fmt.Sprintf("command%d", i),
			Priority: i,
		}
		cr.Commands = append(cr.Commands, cmd)
	}

	return &cr
}

func TestCommandRouterSort(t *testing.T) {
	cases := map[string]struct {
		priorities []int
		expected   []int
	}{
		"already sorted": {
			priorities: []int{1, 2, 3},
			expected:   []int{1, 2, 3},
		},
		"needs sorting": {
			priorities: []int{3, 1, 2},
			expected:   []int{1, 2, 3},
		},
		"same priorities": {
			priorities: []int{0, 0, 0},
			expected:   []int{0, 0, 0},
		},
		"empty list": {
			priorities: []int{},
			expected:   []int{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cr := createCommandRouterFromPriorities(tc.priorities)
			cr.Sort()

			sorted := createCommandRouterFromPriorities(tc.expected)

			assert.Equal(t, sorted.Commands, cr.Commands)
		})
	}
}

func createCommandRouterFromNames(in []string) *CommandRouter {
	cr := CommandRouter{}

	for _, i := range in {
		cmd := &HttpCommand{Command: i}
		cr.Commands = append(cr.Commands, cmd)
	}

	return &cr
}

func TestCommandRouterNames(t *testing.T) {
	cases := map[string]struct {
		commands []string
		expected []string
	}{
		"needs sorting": {
			commands: []string{"def", "abc", "ghi"},
			expected: []string{"abc", "def", "ghi"},
		},
		"already sorted": {
			commands: []string{"abc", "def", "ghi"},
			expected: []string{"abc", "def", "ghi"},
		},
		"no commands": {
			commands: []string{},
			expected: []string{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cr := createCommandRouterFromNames(tc.commands)
			out := cr.Names()

			assert.Equal(t, tc.expected, out)
		})
	}
}

func TestCommandRouterRoute(t *testing.T) {
	cases := map[string]struct {
		commands  []string
		command   string
		returnErr bool
	}{
		"one command": {
			commands:  []string{"command"},
			command:   "command",
			returnErr: false,
		},
		"two commands, found first": {
			commands:  []string{"command1", "command2"},
			command:   "command1",
			returnErr: false,
		},
		"two commands, found second": {
			commands:  []string{"command1", "command2"},
			command:   "command2",
			returnErr: false,
		},
		"no commands": {
			commands:  []string{},
			command:   "command",
			returnErr: true,
		},
		"one command, not found": {
			commands:  []string{"command1"},
			command:   "command2",
			returnErr: true,
		},
		"two commands, not found": {
			commands:  []string{"command1", "command2"},
			command:   "command3",
			returnErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cr := createCommandRouterFromNames(tc.commands)

			out, err := cr.Route(tc.command)

			if !tc.returnErr {
				assert.Equal(t, tc.command, out.GetCommand())
			} else {
				assert.Equal(t, err.Error(), noCommandRoutedErrMsg)
			}
		})
	}
}
