package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

//Stucture of the server
type server struct {
	connected map[net.Addr]*client
	rooms     map[string]*room
	commands  chan command
}

//Function for creating the server
func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

//Function for handeling the client commands
func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMDNick:
			s.nick(cmd.client, cmd.args[1])
		case CMDJoin:
			s.join(cmd.client, cmd.args[1])
		case CMDRooms:
			s.listRooms(cmd.client)
		case CMDMsg:
			s.msg(cmd.client, cmd.args)
		case CMDHelp:
			s.help(cmd.client)
		case CMDQuit:
			s.quit(cmd.client)
		}
	}
}

//Function for creating a client for the connection
func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.msg(fmt.Sprintf("Hello, you have joined the server with the nick of: %s", c.nick))
	c.msg("To see the commands type /help")
	c.readInput()
}

//Function for changing the client nickname
func (s *server) nick(c *client, nick string) {
	c.nick = nick
	c.msg(fmt.Sprintf("Your new nickname has been set to: %s", nick))
}

//Function for joining a room
func (s *server) join(c *client, roomName string) {
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

//Function for listing all rooms
func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

//Function for broadcasting a message to current room
func (s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

//Function for telling the commands of the server
func (s *server) help(c *client) {
	c.msg(fmt.Sprintf("The commands this server supports are \n/nick [new nickname] for setting a new nickname\n/join [room name] for joining or switching rooms\n/rooms for listing all rooms active\n/msg [message] for messaging the current room\n/quit for quitting the server"))
}

//Function for quitting the server
func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("Closing connection...")
	c.conn.Close()
}

//Function for leaving rooms
func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
