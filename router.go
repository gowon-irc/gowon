package main

import "errors"

const (
	noCommandRoutedErrMsg = "no command could be routed"
)

type CommandRouter struct {
	Commands []Command
}

func (cr CommandRouter) Route(command string) (string, error) {
	for _, cmd := range cr.Commands {
		if command == cmd.Command {
			return cmd.Command, nil
		}
	}

	return "", errors.New(noCommandRoutedErrMsg)
}
