package main

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type message struct {
	Msg  string `json:"msg"`
	Nick string `json:"nick,omitempty"`
	Dest string `json:"dest"`
}

func createMessageStruct(body []byte) (m message, err error) {
	err = json.Unmarshal(body, &m)
	if err != nil {
		return m, errors.Wrap(err, "message couldn't be parsed as message json")
	}

	if m.Msg == "" {
		return m, errors.New("message body does not contain any message content")
	}

	if m.Dest == "" {
		return m, errors.New("message body does not contain a destination")
	}

	return m, nil
}

func createMessageBody(dest, msg, nick string) (body []byte, err error) {
	m := &message{
		Dest: dest,
		Msg:  msg,
		Nick: nick,
	}

	body, err = json.Marshal(m)
	if err != nil {
		return body, errors.Unwrap(err)
	}

	return body, nil
}
