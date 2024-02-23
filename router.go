package main

import "errors"

const (
	noCommandRoutedErrMsg = "no command could be routed"
)

type RouterCommand struct {
	Command  string
	Priority int
}

type CommandRouter struct {
	Commands []*RouterCommand
}

func (cr *CommandRouter) Add(cmd *Command) {
	new := &RouterCommand{
		Command:  cmd.Command,
		Priority: cmd.Priority,
	}

	cr.Commands = append(cr.Commands, new)
}

func (cr *CommandRouter) Route(command string) (cmd *RouterCommand, err error) {
	for _, cmd := range cr.Commands {
		if command == cmd.Command {
			return cmd, nil
		}
	}

	return cmd, errors.New(noCommandRoutedErrMsg)
}
