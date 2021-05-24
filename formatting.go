package main

import (
	"fmt"
	"strings"
)

func colourMsg(msg string) string {
	cs := [][]string{
		{"white", "00"},
		{"black", "01"},
		{"blue", "02"},
		{"green", "03"},
		{"red", "04"},
		{"brown", "05"},
		{"magenta", "06"},
		{"orange", "07"},
		{"yellow", "08"},
		{"lgreen", "09"},
		{"cyan", "10"},
		{"lcyan", "11"},
		{"lblue", "12"},
		{"pink", "13"},
		{"grey", "14"},
		{"gray", "14"},
		{"lgrey", "15"},
		{"lgray", "15"},
		{"clear", "99"},
	}

	for _, c := range cs {
		token := fmt.Sprintf("{%s}", c[0])
		code := fmt.Sprintf("\x03%s", c[1])
		msg = strings.ReplaceAll(msg, token, code)
	}

	return msg
}

func chunkMsg(msg string, length int) (chunk, remaining string) {
	// if the msg is shorter than the chunk, no need to split
	if len(msg) <= length {
		return msg, ""
	}

	// find the largest chunk smaller than the given length
	checked := msg
	li := strings.LastIndex(checked, " ")

	for li != -1 {
		if li <= length {
			return strings.TrimSpace(msg[0:li]), strings.TrimSpace(msg[li:])
		}

		checked = msg[0:li]
		li = strings.LastIndex(checked, " ")
	}

	// if the desired condition could not be met, we need to split on the given length
	return msg[0:length], msg[length:]
}

func splitMsg(msg string, length int) (out []string) {
	chunk, remaining := chunkMsg(msg, length)
	out = append(out, chunk)

	for remaining != "" {
		chunk, remaining = chunkMsg(remaining, length)
		out = append(out, chunk)
	}

	return out
}
