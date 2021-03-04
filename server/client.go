package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

//The strcuture of each client
type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

//fuction for listening to the net.Conn and parsing the commands
func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMDNick,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMDJoin,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMDRooms,
				client: c,
			}
		case "/online":
			c.commands <- command{
				id:     CMDOnline,
				client: c,
			}
		case "/msg":
			c.commands <- command{
				id:     CMDMsg,
				client: c,
				args:   args,
			}
		case "/private":
			c.commands <- command{
				id:     CMDPmsg,
				client: c,
				args:   args,
			}
		case "/help":
			c.commands <- command{
				id:     CMDHelp,
				client: c,
			}
		case "/quit":
			c.commands <- command{
				id:     CMDQuit,
				client: c,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

//Functions for sending messages back to client

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
