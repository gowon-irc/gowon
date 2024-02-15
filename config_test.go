package main

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDataDir = filepath.Join("testdata", "config")
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

func TestNewConfigManager(t *testing.T) {
	cm := NewConfigManager()

	assert.Len(t, cm.ConfigFiles, 0)
}

func TestConfigOpenFile(t *testing.T) {
	cases := map[string]struct {
		fn     string
		err    error
		length int
	}{
		"blank config": {
			fn:     "empty.yaml",
			err:    nil,
			length: 1,
		},
		"all required fields": {
			fn:     "required.yaml",
			err:    nil,
			length: 1,
		},
		"invalid yaml": {
			fn:     "invalid.yaml",
			err:    assert.AnError,
			length: 0,
		},
		"invalid fn": {
			fn:     "nonexistent.yaml",
			err:    assert.AnError,
			length: 0,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fn := filepath.Join(testDataDir, tc.fn)

			cm := NewConfigManager()
			err := cm.OpenFile(fn)

			assert.Len(t, cm.ConfigFiles, tc.length)

			if tc.err == nil {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestConfigManagerLoadDirectory(t *testing.T) {
	cases := map[string]struct {
		dir    string
		err    error
		length int
	}{
		"working directory": {
			dir:    "working",
			err:    nil,
			length: 2,
		},
		"empty directory": {
			dir:    "empty",
			err:    nil,
			length: 0,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(testDataDir, "LoadDirectory", tc.dir)

			cm := NewConfigManager()
			err := cm.LoadDirectory(dir)

			assert.Len(t, cm.ConfigFiles, tc.length)
			assert.Nil(t, err)
		})
	}
}

func TestAddOpts(t *testing.T) {
	server := "irc.kwlchat.net:6997"
	cfg := Config{
		Server: server,
	}

	cm := NewConfigManager()
	cm.AddOpts(cfg)

	assert.Equal(t, server, cm.Opts.Server)
}

func TestConfigManagerMerge(t *testing.T) {
	cases := map[string]struct {
		fns      []string
		mergedfn string
		err      error
	}{
		"merge empty with valid": {
			fns:      []string{"empty.yaml", "required.yaml"},
			mergedfn: "empty_required.yaml",
			err:      nil,
		},
		"merge valid with same valid": {
			fns:      []string{"required.yaml", "required.yaml"},
			mergedfn: "empty_required.yaml",
			err:      nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cm := NewConfigManager()
			for _, fn := range tc.fns {
				ffn := filepath.Join(testDataDir, fn)
				_ = cm.OpenFile(ffn)
			}

			cm2 := NewConfigManager()
			_ = cm2.OpenFile(filepath.Join(testDataDir, "merged", tc.mergedfn))
			expected, _ := cm2.Merge()

			out, err := cm.Merge()

			assert.Equal(t, expected, out)
			assert.Nil(t, err)
		})
	}
}
