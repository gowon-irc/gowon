package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessageStructNoError(t *testing.T) {
	body := []byte(`{"msg": "m", "nick": "n", "dest": "d"}`)
	_, err := createMessageStruct(body)

	assert.Nil(t, err)
}

func TestCreateMessageStructErrors(t *testing.T) {
	cases := []struct {
		name   string
		body   []byte
		errMsg string
	}{
		{
			name:   "No content",
			body:   []byte(`{"dest": "d"}`),
			errMsg: "message body does not contain any message content",
		},
		{
			name:   "No destination",
			body:   []byte(`{"msg": "m"}`),
			errMsg: "message body does not contain a destination",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := createMessageStruct(tc.body)
			assert.EqualError(t, err, tc.errMsg)
		})
	}
}
