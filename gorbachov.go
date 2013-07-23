package main

import (
    "fmt"
    "net"
    "log"
    "bufio"
)

type Client struct {
    connection *net.Conn
    buffers *bufio.ReadWriter
    input chan string
    output chan string
    die chan bool
}

func (client *Client) connect(addr string, port uint) {
    connection_string := fmt.Sprintf("%v:%v", addr, port)
    connection, err := net.Dial("tcp", connection_string)
    if err != nil {
        log.Fatalf("Can't connect to server at %v", connection_string)
    }
    client.connection = &connection
    client.buffers = bufio.NewReadWriter(
        bufio.NewReader(connection),
        bufio.NewWriter(connection),
    )
    client.input = make(chan string, 100)
    client.output = make(chan string, 100)
    client.die = make(chan bool)
}

func (client *Client) launch() {

    go func(r *bufio.ReadWriter, readChan chan string) {
        var (
            line []byte
            err error
        )
        for {
            line, err = r.ReadBytes('\r')
            if err != nil {
                log.Fatal("Error while reading from server.")
            }
            readChan <- string(line)
        }
    }(client.buffers, client.output)

    go func(w *bufio.ReadWriter, writeChan chan string) {
        var line string
        for {
            line = <-writeChan
            w.WriteString(line)
            w.Flush()
        }
    }(client.buffers, client.input)

    <-client.die
}

func main() {
    client := &Client{}
    client.connect("irc.freenode.org", 6667)
    input := client.input
    output := client.output
    die := client.die
    go client.launch()

    go func(output chan string) {
        for {
            fmt.Printf("%v", <-output)
        }
    }(output)

    input <- "NICK Gorbachov\n"
    input <- "USER Gorbachov 0 * :Gorbachov\n"
    input <- "JOIN #limajs\n"
    input <- "PRIVMSG #limajs Hai\n"
    for {
        input <- <-output
    }
    die <- true
}
