package message

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type message struct {
	Msg  string `json:"msg"`
	Nick string `json:"nick,omitempty"`
	Dest string `json:"dest"`
}

const ErrorMessageParseMsg = "message couldn't be parsed as message json"
const ErrorMessageNoBodyMsg = "message body does not contain any message content"
const ErrorMessageNoDestinationMsg = "message body does not contain a destination"

func CreateMessageStruct(body []byte) (m message, err error) {
	err = json.Unmarshal(body, &m)
	if err != nil {
		return m, errors.Wrap(err, ErrorMessageParseMsg)
	}

	if m.Msg == "" {
		return m, errors.New(ErrorMessageNoBodyMsg)
	}

	if m.Dest == "" {
		return m, errors.New(ErrorMessageNoDestinationMsg)
	}

	return m, nil
}

func CreateMessageBody(dest, msg, nick string) (body []byte, err error) {
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
