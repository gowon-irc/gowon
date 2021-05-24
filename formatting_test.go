package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColourMsg(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "One colour",
			in:       "{red}Hello world",
			expected: "\x0304Hello world",
		},
		{
			name:     "Colour and clear",
			in:       "{red}Hello{clear} world",
			expected: "\x0304Hello\x0399 world",
		},
		{
			name:     "Two colours and clear",
			in:       "{red}Hello {blue}world{clear}",
			expected: "\x0304Hello \x0302world\x0399",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := colourMsg(tc.in)
			assert.Equal(t, tc.expected, msg)
		})
	}
}

func TestLastColour(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "Multiple colours",
			in:       "\x0302blue\x0399\x0303green",
			expected: "\x0303",
		},
		{
			name:     "One colour",
			in:       "\x0302blue",
			expected: "\x0302",
		},
		{
			name:     "Last colour clear",
			in:       "\x0302blue\x0399\x0303green\x0399",
			expected: "",
		},
		{
			name:     "No colours",
			in:       "Hello",
			expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, lastColour(tc.in))
		})
	}
}

func TestTerminateColour(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "No terminating colour",
			in:       "\x0302Blue string",
			expected: "\x0302Blue string\x0399",
		},
		{
			name:     "Terminating colour",
			in:       "\x0302Blue string\x0399",
			expected: "\x0302Blue string\x0399",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, terminateColour(tc.in))
		})
	}
}

func TestReverseWords(t *testing.T) {
	cases := []struct {
		name     string
		in       []string
		expected []string
	}{
		{
			name:     "Empty string",
			in:       []string{""},
			expected: []string{""},
		},
		{
			name:     "One word",
			in:       []string{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "Multiple words",
			in:       []string{"hello", "world", "goodbye", "world"},
			expected: []string{"world", "goodbye", "world", "hello"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, reverseWords(tc.in))
		})
	}
}

func TestChunkMsg(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		chunks []string
	}{
		{
			name:   "Exact length chunk",
			in:     "I am long, please split me but just once!",
			chunks: []string{"I am long,", "please split me but just once!"},
		},
		{
			name:   "Short message",
			in:     "I am short",
			chunks: []string{"I am short", ""},
		},
		{
			name:   "First word longer",
			in:     "aaaaaaaaaaa first word is long",
			chunks: []string{"aaaaaaaaaa", "a first word is long"},
		},
		{
			name:   "Shorter chunk",
			in:     "I'm long, please split me but just once!",
			chunks: []string{"I'm long,", "please split me but just once!"},
		},
		{
			name:   "No spaces",
			in:     "aaaaaaaaaabb",
			chunks: []string{"aaaaaaaaaa", "bb"},
		},
		{
			name:   "Empty string",
			in:     "",
			chunks: []string{"", ""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chunk, remaining := chunkMsg(tc.in, 10)
			chunks := []string{chunk, remaining}

			assert.Equal(t, tc.chunks, chunks)
		})
	}
}

func TestChunkMsgColour(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		chunks []string
	}{
		{
			name:   "Coloured string",
			in:     "\x0302Both these lines should be blue\x0399",
			chunks: []string{"\x0302Both these lines\x0399", "\x0302should be blue\x0399"},
		},
		{
			name:   "Non terminated coloured string",
			in:     "\x0302The colour does not terminate",
			chunks: []string{"\x0302The colour does not\x0399", "\x0302terminate\x0399"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chunk, remaining := chunkMsg(tc.in, 25)
			chunks := []string{chunk, remaining}

			assert.Equal(t, tc.chunks, chunks)
		})
	}
}

func TestSplitMsg(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		expected []string
	}{
		{
			name:     "Long line",
			in:       "I am a long line please split me",
			expected: []string{"I am a long line", "please split me"},
		},
		{
			name:     "Longer line",
			in:       "I am very long please split me a few times",
			expected: []string{"I am very long", "please split me a", "few times"},
		},
		{
			name:     "Colours",
			in:       "\x0302Blue words and also\x0399 some normal words",
			expected: []string{"\x0302Blue words and\x0399", "\x0302also\x0399 some", "normal words"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msgs := splitMsg(tc.in, 20)
			assert.Equal(t, tc.expected, msgs)
		})
	}
}
