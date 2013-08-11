package gorbachov

import (
	"errors"
	"fmt"
	"testing"
)

type parserOutput struct {
	prefix  string
	command string
	args    []string
}

func (po *parserOutput) checkOutput(prefix, command string, args []string) error {
	if po.prefix != prefix {
		return errors.New("Prefix doesn't match.")
	}

	if po.command != command {
		return errors.New("Command doesn't match.")
	}

	if len(po.args) != len(args) {
		for _, v := range args {
			fmt.Printf("'%v'\n", v)
		}
		err := fmt.Sprintf("%v(%v) != %v(%v)", po.args, len(po.args), args, len(args))
		return errors.New("Argument list size doesn't match: " + err)
	}

	for index, value := range po.args {
		if value != args[index] {
			err := fmt.Sprintf("%v != %v", value, args[index])
			return errors.New("Arguments don't match." + err)
		}
	}
	return nil
}

func TestParser(t *testing.T) {

	// Test various IRC messages
	inputs := make(map[string]*parserOutput)
	inputs["PING :kornbluth.freenode.net"] = &parserOutput{
		"",
		"PING",
		[]string{"kornbluth.freenode.net"},
	}
	inputs[":BotName MODE BotName :+i"] = &parserOutput{
		"BotName",
		"MODE",
		[]string{"BotName", "+i"},
	}
	inputs[":someone!~someone@someserver NOTICE BotName :PING 137498005"] = &parserOutput{
		"someone!~someone@someserver",
		"NOTICE",
		[]string{"BotName", "PING 137498005"},
	}
	inputs[":someone!~someone@someserver PRIVMSG #somechannel :Hello my friend"] = &parserOutput{
		"someone!~someone@someserver",
		"PRIVMSG",
		[]string{"#somechannel", "Hello my friend"},
	}
	inputs[":someone!~someone@someserver PRIVMSG BotName :sup"] = &parserOutput{
		"someone!~someone@someserver",
		"PRIVMSG",
		[]string{"BotName", "sup"},
	}

	for input, out_struct := range inputs {
		prefix, command, arguments, _ := ParseMessage(input)
		if err := out_struct.checkOutput(prefix, command, arguments); err != nil {
			t.Error(err)
		}
	}

	// test blank messages
	if _, _, _, err := ParseMessage(""); err == nil {
		t.Error("Blank message should throw error.")
	}
	if _, _, _, err := ParseMessage(" "); err == nil {
		t.Error("Blank message should throw error.")
	}
	if _, _, _, err := ParseMessage("  "); err == nil {
		t.Error("Blank message should throw error.")
	}

}
