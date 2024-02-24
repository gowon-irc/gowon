package main

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gowon-irc/go-gowon"
)

const (
	noCommandRoutedErrMsg = "no command could be routed"
)

type RouterCommand interface {
	Send(in *gowon.Message) *gowon.Message
	GetCommand() string
	GetPriority() int
	Match(string) bool
}

type HttpCommand struct {
	Command  string
	Regex    string
	Priority int
}

func (hc *HttpCommand) Send(in *gowon.Message) (out *gowon.Message) {
	return &gowon.Message{
		Module:  hc.Command,
		Msg:     fmt.Sprintf("message from %s", hc.Command),
		Nick:    "tester",
		Dest:    "#gowon",
		Command: ".test",
		Args:    "command",
	}
}

func (hc *HttpCommand) GetCommand() string {
	return hc.Command
}

func (hc *HttpCommand) GetPriority() int {
	return hc.Priority
}

func (hc *HttpCommand) Match(text string) bool {
	if hc.Command == gowon.GetCommand(text) {
		return true
	}

	if hc.Regex != "" {
		re := regexp.MustCompile(hc.Regex)
		return re.Match([]byte(text))
	}

	return false
}

type InternalCommand struct {
	Command  string
	Priority int
	f        func() string
}

func (ic *InternalCommand) Send(in *gowon.Message) (out *gowon.Message) {
	return &gowon.Message{
		Module:  ic.Command,
		Msg:     ic.f(),
		Nick:    in.Nick,
		Dest:    in.Dest,
		Command: ic.Command,
		Args:    "",
	}
}

func (ic *InternalCommand) GetCommand() string {
	return ic.Command
}

func (ic *InternalCommand) GetPriority() int {
	return ic.Priority
}

func (ic *InternalCommand) Match(text string) bool {
	return ic.Command == gowon.GetCommand(text)
}

type CommandRouter struct {
	Commands []RouterCommand
}

func (cr *CommandRouter) Add(cmd *Command) {
	new := &HttpCommand{
		Command:  cmd.Command,
		Priority: cmd.Priority,
		Regex:    cmd.Regex,
	}

	cr.Commands = append(cr.Commands, new)
}

func (cr *CommandRouter) AddInternal(command string, f func() string) {
	new := &InternalCommand{
		Command:  command,
		Priority: -99,
		f:        f,
	}

	cr.Commands = append(cr.Commands, new)
}

func (cr *CommandRouter) Sort() {
	sort.Slice(cr.Commands, func(i, j int) bool {
		return cr.Commands[i].GetPriority() < cr.Commands[j].GetPriority()
	})
}

func (cr *CommandRouter) Names() []string {
	out := []string{}

	for _, c := range cr.Commands {
		out = append(out, c.GetCommand())
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})

	return out
}

func (cr *CommandRouter) Route(text string) (RouterCommand, error) {
	for _, cmd := range cr.Commands {
		if cmd.Match(text) {
			return cmd, nil
		}
	}

	return nil, errors.New(noCommandRoutedErrMsg)
}

func colourList(in []string) (out []string) {
	out = []string{}

	colours := []string{"green", "red", "blue", "orange", "magenta", "cyan", "yellow"}
	cl := len(colours)

	for n, i := range in {
		c := colours[n%cl]
		o := fmt.Sprintf("{%s}%s{clear}", c, i)
		out = append(out, o)
	}

	return out
}

func createHelpCommandFunc(cr *CommandRouter) func() string {
	return func() string {
		return strings.Join(colourList(cr.Names()), ", ")
	}
}
