package gorbachov

import (
	"errors"
	"strings"
)

// ParseMessage splits a raw IRC message into it's prefix, command and argument list
// returns an error if raw string was blank.
// Adapted from https://github.com/powdahound/twisted/blob/master/twisted/words/protocols/irc.py#L75
func ParseMessage(raw string) (prefix, command string, arguments []string, err error) {
	prefix, command = "", ""
	arguments, temp_args := []string{}, []string{}
	err = nil

	if len(raw) == 0 || strings.Replace(raw, " ", "", -1) == "" {
		err = errors.New("Blank message")
		return
	}

	if raw[0] == ':' {
		splitted := strings.SplitN(raw[1:], " ", 2)
		prefix, raw = splitted[0], splitted[1]
	}
	if strings.Index(raw, " :") != -1 {
		trailing := ""
		splitted := strings.SplitN(raw, ":", 2)
		raw, trailing = splitted[0], splitted[1]
		temp_args = strings.Split(raw, " ")
		temp_args = append(temp_args, trailing)
	} else {
		temp_args = strings.Split(raw, " ")
	}

	for _, v := range temp_args {
		if v != " " && len(v) != 0 {
			arguments = append(arguments, v)
		}
	}
	command = arguments[0]
	arguments = arguments[1:]
	return
}
