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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chunk, remaining := chunkMsg(tc.in, 10)
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msgs := splitMsg(tc.in, 20)
			assert.Equal(t, tc.expected, msgs)
		})
	}
}
