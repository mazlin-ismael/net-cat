package server

import (
	"errors"
	"net"
	"os"
	"reflect"
)

// get the penguin message
func readPenguinFile() []byte {
	penguinTXT, errFile := os.ReadFile("SERVER/TXT/penguinMSG.txt")
	if errFile != nil {
		return nil
	}
	return penguinTXT
}

// init a name of a user
func (s *Server) initUserName(conn net.Conn, msg string) ([]byte, error) {
	var userName []byte
	for {
		var restartLoop bool
		conn.Write([]byte(msg))
		bufName := make([]byte, 2048)
		n, err := conn.Read(bufName)
		if err != nil {
			return nil, errors.New("error conn read")
		}

		if len(bufName[:n]) <= 1 {
			conn.Write([]byte("Only accept connection with non-empty name\n\n"))
			continue
		}

		for connClient, name := range s.Clients {
			if reflect.DeepEqual(name.Username, bufName[:n-1]) && connClient != conn {
				conn.Write([]byte("Pseudo already taken\n\n"))
				restartLoop = true
			}
		}
		if restartLoop {
			continue
		}
		userName = bufName[:n-1]
		break
	}
	return userName, nil
}

// transmetting the message to the main
func (s *Server) TransmettingMSG(message Message) {
	s.Msgch <- message
}

// init the message
func (s *Server) InitMSG(conn net.Conn, typeMSG string, content []byte) Message {
	var msg Message = Message{
		From:    s.Clients[conn],
		TypeMSG: typeMSG,
	}
	// define the content of the message
	msg.definePayload(content)
	return msg
}
