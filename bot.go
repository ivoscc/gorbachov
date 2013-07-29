package main

import (
    "log"
    "fmt"
    "regexp"
)

type Bot struct {
    server string
    name string
    channel string
    message_handlers map[string]func(string) string
    client *Client
}

func (bot *Bot) say(message string) {
    message = fmt.Sprintf("PRIVMSG %v :%v", bot.channel, message)
    bot.client.SendResponse(message)
}

// Dispatches a handler function depending on the PRIVMSG received
func (bot *Bot) HandleMessage(message string) {
    for str, handler := range bot.message_handlers {
        matches, err := regexp.MatchString(str, message)

        if err != nil {
            log.Fatal(err)
        }

        if matches {
            response := (handler)(message)
            bot.say(response)
        }
    }
}

func (bot *Bot) start() {
    client := bot.client
    client.Connect(bot.server)
    message_stream := client.Start(bot.name, bot.channel)

    for {
        message := <-message_stream
        arguments := message.arguments
        bot.HandleMessage(arguments[len(arguments)-1])
    }
}

// test echo bot
func main() {
    mapping := make(map[string]func(string) string)
    // catch all function
    mapping[".*"] = func (message string) string {
        return message
    }

    bot := &Bot{
        "irc.freenode.org:6667",
        "Gorbachov",
        "#chorbagov",
        mapping,
        &Client{},
    }
    bot.start()
}
