package message

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type Message struct {
	Module  string `json:"module"`
	Msg     string `json:"msg"`
	Nick    string `json:"nick,omitempty"`
	Dest    string `json:"dest"`
	Command string `json:"command"`
	Args    string `json:"args"`
}

const ErrorMessageParseMsg = "message couldn't be parsed as message json"

const ErrorMessageNoModuleMsg = "message body does not contain a module source"

const ErrorMessageNoBodyMsg = "message body does not contain any message content"

const ErrorMessageNoDestinationMsg = "message body does not contain a destination"

func getCommand(msg string) string {
	if strings.HasPrefix(msg, ".") {
		return strings.TrimPrefix(strings.Fields(msg)[0], ".")
	}

	return ""
}

func getArgs(msg string) string {
	if !strings.HasPrefix(msg, ".") {
		return msg
	}

	return strings.TrimSpace(strings.TrimPrefix(msg, strings.Fields(msg)[0]))
}

func CreateMessageStruct(body []byte) (m Message, err error) {
	err = json.Unmarshal(body, &m)
	if err != nil {
		return m, errors.Wrap(err, ErrorMessageParseMsg)
	}

	if m.Module == "" {
		return m, errors.New(ErrorMessageNoModuleMsg)
	}

	if m.Msg == "" {
		return m, errors.New(ErrorMessageNoBodyMsg)
	}

	if m.Dest == "" {
		return m, errors.New(ErrorMessageNoDestinationMsg)
	}

	if m.Command == "" {
		m.Command = getCommand(m.Msg)
	}

	if m.Args == "" {
		m.Args = getArgs(m.Msg)
	}

	return m, nil
}
