package main

import (
	"net"
	serverNC "NETCAT/SERVER"
)

// list of global vars used in is project
var (
	defaultPort   string = "8988"
	port          string
	serverRank    int
	currentView             = 0
	currentServer           = 0
	portWaiting   chan bool = make(chan bool)
	initBinds     chan bool = make(chan bool)
	firstTime     bool      = true
	portsUsed     []string
)

// struct use for the UI interface
type ViewsSettings struct {
	Users []struct 
	{
		Conn 		net.Conn
		Name 		string
		ViewName 	string
	}
	server serverNC.Server
}