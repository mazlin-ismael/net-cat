package server

import (
	"net"
	"fmt"
)

// init new Server
func NewServer(listenAddr string, rank int) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Rank: 		rank,
		PenguinMSG: readPenguinFile(),
		Quitch:    	make(chan bool),
		Msgch: 		make(chan Message),
		Clients: 	make(map[net.Conn]*Client),
	}
}

// start new Server
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.Ln = ln
	go s.acceptLoop()
	<- s.Quitch
	return nil
}


// listen the incoming connections
func (s *Server) acceptLoop() {
	for {
		// more that 10 connections, refuse the others conn
		if len(s.Clients) < 10 {
			conn, err := s.Ln.Accept()
			if err != nil {
				fmt.Println("accept error:", err)
				continue
			}
			// init the client
			s.Clients[conn] = &Client{
				Conn: 		conn,
				ViewName: conn.RemoteAddr().String(),
			}
			// transmetting a msg to make a new View of the conn
			s.TransmettingMSG(Message{
				From: s.Clients[conn],
				TypeMSG: "newView",
			})
			// for each connexion
			go s.readLoop(conn)
		}
	}
}