package serverNC

import (
	"fmt"
	"net"
	"os"
)

type Message struct {
	From []byte
	Payload []byte  
}

type Server struct {
	ListenAddr 	string
	Ln 			net.Listener
	Quitch		chan struct{}
	Msgch       chan Message
	EntryMsg	[]byte
	Clients 	map[string]net.Conn
	// PeerMap map[net.Addr]
}

func NewServer(listenAddr string) *Server {

	entryMsg, errMsg := os.ReadFile("server/penguinMSG.txt")
	if errMsg != nil {
		fmt.Println("error file entry message:", entryMsg)
	}

	return &Server {
		ListenAddr: listenAddr,
		Quitch: 	make(chan struct{}),
		Msgch: 		make(chan Message, 10),
		EntryMsg:	entryMsg,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil { 
		return err
	}

	defer ln.Close()
	s.Ln = ln

	go s.acceptLoop()

	<- s.Quitch
	close(s.Msgch)
	return nil
}

func (s *Server) acceptLoop() {
	var count int
	for {
		if count == 3 {
			s.Ln.Close()
		} else {
			count++
			conn, err := s.Ln.Accept()
			if err != nil {
				fmt.Println("accept error:", err)
				continue 
			}
			fmt.Println("new connection to the server:", conn.RemoteAddr() )
			go s.readLoop(conn)
		}

	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	conn.Write(s.EntryMsg)
	var userName []byte
	for {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		bufName := make([]byte, 2048)
		n, err := conn.Read(bufName)
		if err != nil {
			fmt.Println("read error", err)
			return
		}
		if len(bufName[:n]) <= 1 {
			conn.Write([]byte("Only accept connection with non-empty name\n\n"))
			continue
		}
		userName = bufName[:n]
		break
	}
	
	for _, client := range s.Clients {
		client.Conn.Write([]byte(string(userName[:len(userName)-1])+" has joined our chat...\n"))
	}
	s.Clients[string(userName)] = conn

	precdentsMSG, _ := os.ReadFile("logs.txt")
	conn.Write([]byte(precdentsMSG))
	

	for {
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			for i, client := range s.Clients {
				if client.Conn == conn {
					s.Clients = append(s.Clients[0:i], s.Clients[i+1:]...)
				} else {
					client.Conn.Write([]byte(string(userName[:len(userName)-1])+"has left our chat...\n"))
				}
			}
			return
		}

		if len(buf[:n]) <= 1 {
			conn.Write([]byte("Only accept non-empty message\n"))
			continue
		}
		fmt.Println(userName)
		s.Msgch <- Message{
			From: 		userName,
			Payload:	buf[:n],
		}
	} 
}