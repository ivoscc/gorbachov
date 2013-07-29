package main

import (
    "fmt"
    "log"
    "net/textproto"
)

// The IRC client
type Client struct {
    connection *textproto.Conn
}

// Server message (http://tools.ietf.org/html/rfc1459.html#section-2.3)
type Message struct {
    prefix string
    command string
    arguments []string
}

// SendResponse sends a 'message' string to the client's connection.
func (client *Client) SendResponse(message string) {
    log.Printf("=> %v", message)
    connection := client.connection
    _, err := connection.Cmd(message)
    if err != nil {
        log.Fatal(err)
    }
}

// ReadLine reads a line from the open connection and prints it.
func (client *Client) ReadLine() string {
    line, err := client.connection.ReadLine()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("<= %v\n", line)
    return line
}

// Connect opens a connection to an IRC server
func (client *Client) Connect(server string) {

	conn, err := textproto.Dial("tcp", server)
	if err != nil {
        log.Fatal(err)
	}
    client.connection = conn
}

// Start dispatches the client and returns a channel to receive messages from
func (client *Client) Start(botname, channel string) chan Message {

    established := make(chan bool)

    // TODO: timeout
    go func(established chan bool) {
        for {
            _, command, _, err := ParseMessage(client.ReadLine())
            if err != nil {
                log.Fatalf("Error while parsing received message")
            }
            if command == "MODE" {
                break
            }
        }
        established <- true
    }(established)

    client.SendResponse(fmt.Sprintf("NICK %v", botname))
    client.SendResponse(fmt.Sprintf("USER %v 0 * %v", botname, botname))
    client.JoinChannel(channel)

    <-established

    message_stream := make(chan Message)
    go client.Run(message_stream)
    return message_stream
}

func (client *Client) Run(message_stream chan Message) {

    for {
        prefix, command, arguments, err := ParseMessage(client.ReadLine())
        if err != nil {
            log.Fatalf("Error while parsing received message")
        }
        if command == "PRIVMSG" {
            message_stream <- Message{prefix, command, arguments}
        } else if command == "PING" {
            client.Execute(fmt.Sprintf("PONG :%v", arguments[0]))
        }
    }
}

func (client *Client) JoinChannel(channel string) {
    // TODO: Validate channel (#, &, etc)
    client.SendResponse(fmt.Sprintf("JOIN %v", channel))
}

// Execute sends an arbitrary command to the server (no validity checking is done)
func (client *Client) Execute(command string) {
    client.SendResponse(command)
}
