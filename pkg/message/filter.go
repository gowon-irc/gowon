package message

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const ErrorFilterInvalidChars = "filter contains invalid characters"

const ErrorFilterInvertNotAtBeginning = "Inversion (!) can only exist at the beginning of the filter"

const ErrorFilterInvalidFields = "Filter does not contain two fields, equals should appear once"

const ErrorFilterInvalidKey = "Filter key is invalid, must be one of module, msg, nick, dest, command, args"

func checkFilterSingle(filter string) error {
	m, err := regexp.MatchString(`^[!=a-zA-Z0-9]+$`, filter)
	if err != nil {
		return err
	}
	if !m {
		return errors.New(ErrorFilterInvalidChars)
	}

	m, err = regexp.MatchString(`[=a-zA-Z]!`, filter)
	if err != nil {
		return err
	}
	if m {
		return errors.New(ErrorFilterInvertNotAtBeginning)
	}

	if strings.Count(filter, "=") != 1 {
		return errors.New(ErrorFilterInvalidFields)
	}

	contains := func(ss []string, c string) bool {
		for _, s := range ss {
			if s == c {
				return true
			}
		}
		return false
	}

	key := strings.TrimPrefix(strings.Split(filter, "=")[0], "!")
	if !contains([]string{"module", "msg", "nick", "dest", "command", "args"}, key) {
		return errors.New(ErrorFilterInvalidKey)
	}

	return nil
}

func CheckFilter(filter string) error {
	for _, f := range strings.Split(filter, "+") {
		err := checkFilterSingle(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func invert(b, inv bool) bool {
	if inv {
		return !b
	}

	return b
}

func filterSingle(m *Message, filter string) (bool, error) {
	err := checkFilterSingle(filter)
	if err != nil {
		return false, err
	}

	out := false

	messageFields := map[string]string{
		"module":  m.Module,
		"msg":     m.Msg,
		"nick":    m.Nick,
		"dest":    m.Dest,
		"command": m.Command,
		"args":    m.Args,
	}

	invertFilter := strings.HasPrefix(filter, "!")
	kv := strings.Split(strings.TrimPrefix(filter, "!"), "=")
	key, value := kv[0], kv[1]

	if messageFields[key] == value {
		out = true
	}

	return invert(out, invertFilter), nil
}

func Filter(m *Message, filter string) (bool, error) {
	for _, f := range strings.Split(filter, "+") {
		res, err := filterSingle(m, f)

		if err != nil {
			return false, err
		}

		if !res {
			return false, nil
		}
	}

	return true, nil
}
