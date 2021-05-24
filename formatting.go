package main

import (
	"fmt"
	"regexp"
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

func lastColour(msg string) (colour string) {
	r, _ := regexp.Compile("\x03\\d{2}")
	colours := r.FindAllString(msg, -1)

	if len(colours) == 0 {
		return ""
	}

	c := colours[len(colours)-1]

	if c == "\x0399" {
		return ""
	}

	return c
}

func terminateColour(msg string) string {
	if lastColour(msg) != "" {
		return msg + "\x0399"
	}

	return msg
}

func reverseWords(words []string) []string {
	for i := 0; i < len(words)/2; i++ {
		j := len(words) - i - 1
		words[i], words[j] = words[j], words[i]
	}

	return words
}

func chunkMsg(msg string, length int) (chunk, remaining string) {
	msg = terminateColour(msg)
	if len(msg) <= length {
		return msg, ""
	}

	for _, word := range reverseWords(strings.Fields(msg)) {
		remaining = fmt.Sprintf(" %s%s", word, remaining)
		chunk = strings.TrimSuffix(msg, remaining)

		c := terminateColour(chunk)
		if len(c) <= length {
			lc := lastColour(chunk)
			remaining = fmt.Sprintf("%s%s", lc, strings.TrimSpace(remaining))
			return c, remaining
		}
	}

	terminated := terminateColour(msg)
	return terminated[0:length], terminated[length:]
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
