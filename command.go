package main

type commandID int

//Server commands
const (
	CMDNick commandID = iota
	CMDJoin
	CMDRooms
	CMDMsg
	CMDHelp
	CMDQuit
)

//The structure of the each command being sent over the channel
type command struct {
	id     commandID
	client *client
	args   []string
}