package main

import (
	"net"
)

//Stcuture of each room
type room struct {
	name    string
	members map[net.Addr]*client
}

//Function for sending message to each member of the room except the sender
func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}