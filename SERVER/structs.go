package server

import "net"

// struct for a message to send
type Message struct {
	From		*Client
	Payload		[]byte
	TypeMSG		string
}

// struct for each Client
type Client struct {
	Username	[]byte
	OldName		[]byte
	Conn 	 	net.Conn
	ViewName    string
}

// struct for the server
type Server struct {
	ListenAddr 	string
	Rank 		int
	PenguinMSG 	[]byte
	Quitch  	chan bool
	Ln 			net.Listener
	Msgch 		chan Message
	Clients 	map[net.Conn]*Client
}