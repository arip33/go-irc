package irc

import (
	"bufio"

	"errors"
	"log"
	"net"
	"strings"
)

type Conn struct {
	conn     *net.TCPConn
	Received chan string
	ToSend   chan string
}

func Dial(server string) (*Conn, error) {
	ipAddr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, ipAddr)
	if err != nil {
		return nil, err
	}

	r := make(chan string, 200)
	w := make(chan string, 200)
	c := &Conn{conn: conn, Received: r, ToSend: w}

	// Reading task
	go func() {
		r := bufio.NewReader(conn)
		for {
			data, err := r.ReadString('\n')
			if err != nil {
				log.Println("Read error: ", err)
				return
			}
			if strings.HasPrefix(data, "PING") {
				c.ToSend <- "PONG" + data[4:len(data)-2]
			} else {
				c.Received <- data[0 : len(data)-2]
			}
		}
	}()

	// Writing task
	go func() {
		w := bufio.NewWriter(conn)
		for {
			data, ok := <-c.ToSend
			if !ok {
				return
			}
			_, err := w.WriteString(data + "\r\n")
			if err != nil {
				log.Println("Write error: ", err)
			}
			w.Flush()
		}
	}()

	return c, nil
}

func (c *Conn) Close() {
}

func (c *Conn) Write(data string) error {
	c.ToSend <- data
	return nil
}

func (c *Conn) Read() (string, error) {
	// blocks until message is available
	data, ok := <-c.Received
	if !ok {
		return "", errors.New("Read stream closed")
	}
	return data, nil
}
