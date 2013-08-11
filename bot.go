package gorbachov

import (
	"fmt"
	"log"
	"regexp"
)

// Server message (http://tools.ietf.org/html/rfc1459.html#section-2.3)
type Message struct {
	Prefix    string
	Command   string
	Arguments []string
}

// GetPrivMSG is a utility function to get the text portion of a Message.
func (message *Message) GetPrivMSG() string {
	return message.Arguments[len(message.Arguments)-1]
}

type Bot struct {
	server   string
	name     string
	channel  string
	handlers map[string]func(*Bot, Message)
	client   *Client
}

// CreateBot initializes a new Bot and returns a pointer to it.
func CreateBot(server, name, channel string) *Bot {
	handlers := make(map[string]func(*Bot, Message))
	return &Bot{server, name, channel, handlers, &Client{}}
}

func (bot *Bot) AddHandler(regex string, handler func(*Bot, Message)) {
	bot.handlers[regex] = handler
}

func (bot *Bot) Say(message string) {
	message = fmt.Sprintf("PRIVMSG %v :%v", bot.channel, message)
	bot.client.SendResponse(message)
}

// Dispatches a handler function depending on the PRIVMSG received
func (bot *Bot) HandleMessage(message Message) {
	arguments := message.Arguments
	for str, handler := range bot.handlers {
		matches, err := regexp.MatchString(str, arguments[len(arguments)-1])

		if err != nil {
			log.Fatal(err)
		}

		if matches {
			go (handler)(bot, message)
		}
	}
}

func (bot *Bot) Start() {
	client := bot.client
	client.Connect(bot.server)
	message_stream := client.Start(bot.name, bot.channel)

	for {
		message := <-message_stream
		bot.HandleMessage(message)
	}
}
