package main

import (
	"bufio"
    "bytes"
	"flag"
	"fmt"
	"irc"
	"os"
	"strings"
)

var server *string = flag.String("server", "irc.freenode.net", "IRC server address")
var port *int = flag.Int("port", 6667, "IRC server port")
var nick *string = flag.String("nick", "go-irc-client", "Nickname")

var help = `
********************************************************************************

JOIN #<name> 					   - join channel
PRIVMSG #<channel name> :<message> - send message to given channel


More info: http://tools.ietf.org/html/rfc1459

********************************************************************************
`

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%v", *server, *port)
	c, err := irc.Dial(addr)
	if err != nil {
		panic(err)
	}

    fmt.Printf("\n** For more information type `help` **\n\n")

	defer c.Close()

	c.Write("NICK " + *nick)
	c.Write("USER bot * * :...")

	// irc messages reader
	go func() {
		for {
			msg, err := c.Read()
			if err != nil {
				panic(fmt.Sprintf("client read error: %s", err))
			}
            if bytes.Equal(msg.Command, []byte("PRIVMSG")) {
                fmt.Printf("%s:: %s -> %s\n", msg.Command, msg.Params, msg.Trailing)
            } else {
                fmt.Println("> ", msg.String())
            }
		}
	}()

	// user input reader
	in := bufio.NewReader(os.Stdin)
	for {
		data, err := in.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("client write error: %s", err))
		}
		data = strings.TrimSpace(data)
		if data == "help" {
			fmt.Println(help)
		} else {
			c.Write(data)
		}
	}
}
