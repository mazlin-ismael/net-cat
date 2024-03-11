package server

import (
	"net"
	"os"
	"strings"
)

var (
	rename   = "/rename"
	users    = "-users"
	commands = "--help"
)

// for each conn
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	// welcome msg of penguin
	conn.Write(s.PenguinMSG)

	// request to the conn a name
	userName, errName := s.initUserName(conn, "[ENTER YOUR NAME]: ")
	// if err delete the user
	if errName != nil {
		s.TransmettingMSG(s.InitMSG(conn, "exit", nil))
		delete(s.Clients, conn)
		return
	}
	_, Noexit := s.Clients[conn]
	if !Noexit {
		return
	}
	// write the precedents logs
	logs, _ := os.ReadFile("LOGS/" + s.ListenAddr + ".log")
	conn.Write([]byte(logs))

	s.Clients[conn].Username = userName

	// announce to others a new user
	s.TransmettingMSG(s.InitMSG(conn, "entrance", nil))

	// listen the messages of user
	s.listenClient(conn)
}

// listen the messages of user
func (s *Server) listenClient(conn net.Conn) {
	for {
		buf := make([]byte, 2048)
		// listen the messages of the user
		n, err := conn.Read(buf)
		// if err delete the user
		if err != nil {
			s.TransmettingMSG(s.InitMSG(conn, "exit", nil))
			delete(s.Clients, conn)
			return
		}
		_, notKicked := s.Clients[conn]
		if !notKicked {
			return
		}
		switch strings.TrimSpace(string(buf[:n])) {
		// if flag rename to rename the name of the user
		case rename:
			userName, errName := s.initUserName(conn, "[ENTER YOUR NEW NAME]: ")
			if errName != nil {
				s.TransmettingMSG(s.InitMSG(conn, "exit", nil))
				delete(s.Clients, conn)
				return
			}
			s.Clients[conn].OldName, s.Clients[conn].Username = s.Clients[conn].Username, userName
			s.TransmettingMSG(s.InitMSG(conn, "rename", nil))
			continue
		// if flag users, list the users with a declared name
		case users:
			for _, client := range s.Clients {
				if client.Username != nil {
					conn.Write([]byte(string(client.Username) + "\n"))
				}
			}
		// if flag help, list commands available
		case commands:
			commands, _ := os.ReadFile("SERVER/TXT/commands.txt")
			conn.Write([]byte(commands))
		// transmetting the message to the others users
		default:
			s.TransmettingMSG(s.InitMSG(conn, "client", buf[:n-1]))
		}
	}
}
