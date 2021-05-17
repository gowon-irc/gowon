package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessageStructNoError(t *testing.T) {
	body := []byte(`{"msg": "m", "nick": "n", "dest": "d"}`)
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
			name:   "No content",
			body:   []byte(`{"dest": "d"}`),
			errMsg: ErrorMessageNoBodyMsg,
		},
		{
			name:   "No destination",
			body:   []byte(`{"msg": "m"}`),
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
