Gorbachov
=========

Mini framework for writing IRC bots in Go.
Just for fun (there are like [1000 others](http://go-lang.cat-v.org/pure-go-libs))

Installation
------------

```bash
go get github.com/ivoscc/gorbachov
```

Usage
-----

To build a simple (and annoying) echo bot that will repeat everything you say, you can do the following.

```go
package main

import (
    "github.com/ivoscc/gorbachov"
)

func EchoHandler(bot *gorbachov.Bot, message gorbachov.Message) {
    text := message.GetPrivMSG()
    bot.Say(text)
}

func main() {
    // Declare the bot
	bot := gorbachov.CreateBot("irc.freenode.org:6667", "Gorbachov", "#some_channel")
    bot.AddHandler(".*", EchoHandler) // You may use any regexp here
	bot.Start()
}
```

And that's it.
