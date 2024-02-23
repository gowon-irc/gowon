package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createCommandRouterFromList(in []string) *CommandRouter {
	cr := CommandRouter{}

	for _, i := range in {
		cmd := &RouterCommand{Command: i}
		cr.Commands = append(cr.Commands, cmd)
	}

	return &cr
}

func TestConfigRoute(t *testing.T) {
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
			cfg := createCommandRouterFromList(tc.commands)

			out, err := cfg.Route(tc.command)

			if !tc.returnErr {
				assert.Equal(t, tc.command, out.Command)
			} else {
				assert.Equal(t, err.Error(), noCommandRoutedErrMsg)
			}
		})
	}
}

func TestCommandRouterAdd(t *testing.T) {
	cr := CommandRouter{}
	cmd := &Command{Command: "command"}
	cr.Add(cmd)

	assert.Equal(t, "command", cr.Commands[0].Command)
}
