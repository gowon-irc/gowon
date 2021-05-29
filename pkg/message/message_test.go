package message

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessageStructNoError(t *testing.T) {
	body := []byte(`{"module": "test", "msg": "m", "nick": "n", "dest": "d"}`)
	_, err := CreateMessageStruct(body)

	assert.Nil(t, err)
}

func TestCreateMessageStructErrors(t *testing.T) {
	cases := []struct {
		name   string
		body   []byte
		errMsg string
	}{
		{
			name:   "No module",
			body:   []byte(`{}`),
			errMsg: ErrorMessageNoModuleMsg,
		},
		{
			name:   "No content",
			body:   []byte(`{"module": "test", "dest": "d"}`),
			errMsg: ErrorMessageNoBodyMsg,
		},
		{
			name:   "No destination",
			body:   []byte(`{"module": "test", "msg": "m"}`),
			errMsg: ErrorMessageNoDestinationMsg,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CreateMessageStruct(tc.body)
			assert.EqualError(t, err, tc.errMsg)
		})
	}
}

func TestMessageCommandArgs(t *testing.T) {
	bodyTmpl := `{"module": "m", "msg": "%s", "dest": "d"}`

	cases := []struct {
		name    string
		body    string
		command string
		args    string
	}{
		{
			name:    "Command with args",
			body:    fmt.Sprintf(bodyTmpl, ".command args"),
			command: "command",
			args:    "args",
		},
		{
			name:    "No command",
			body:    fmt.Sprintf(bodyTmpl, "command args"),
			command: "",
			args:    "command args",
		},
		{
			name:    "Args without command",
			body:    fmt.Sprintf(bodyTmpl, "args without command"),
			command: "",
			args:    "args without command",
		},
		{
			name:    "Existing command and args",
			body:    `{"module": "m", "msg": "output", "dest": "d", "command": "c", "args": "a"}`,
			command: "c",
			args:    "a",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := CreateMessageStruct([]byte(tc.body))

			assert.Nil(t, err)

			got := []string{m.Command, m.Args}
			expected := []string{tc.command, tc.args}

			assert.Equal(t, got, expected)
		})
	}
}
