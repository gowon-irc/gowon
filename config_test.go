package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigServerRegex(t *testing.T) {
	cases := map[string]struct {
		input string
		match bool
	}{
		"server and port": {
			input: "irc.kwlchat.net:6697",
			match: true,
		},
		"ip and port": {
			input: "127.0.0.1:6697",
			match: true,
		},
		"too many port digits": {
			input: "irc.kwlchat.net:999999",
			match: false,
		},
		"port too high": {
			input: "irc.kwlchat.net:65536",
			match: false,
		},
		"zero port": {
			input: "irc.kwlchat.net:0",
			match: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			matched, err := regexp.MatchString(serverRegex, tc.input)

			assert.Equal(t, matched, tc.match)
			assert.Nil(t, err)
		})
	}
}

func TestConfigIrcChannelRegex(t *testing.T) {
	cases := map[string]struct {
		input string
		match bool
	}{
		"channel": {
			input: "#chat",
			match: true,
		},
		"channel with number": {
			input: "#2chat2",
			match: true,
		},
		"channel starting with &": {
			input: "&chat",
			match: true,
		},
		"no hash": {
			input: "chat",
			match: false,
		},
		"channel with a comma": {
			input: "#chat,chat",
			match: false,
		},
		"channel with ^G": {
			input: "#chat\a",
			match: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			matched, err := regexp.MatchString(ircChannelRegex, tc.input)

			assert.Equal(t, matched, tc.match)
			assert.Nil(t, err)
		})
	}
}
