package main

import (
	"fmt"
	"os"
	// serverNC "NETCAT/SERVER"
	// "time"
)

// define a new port to launch the next Server
func listenStdout(nextPort *string) {
	if firstTime {
		if len(os.Args) == 1 {
			newPort(&defaultPort)
			port = defaultPort
		} else if len(os.Args) == 2 {
			port = os.Args[1]
		} else {
			fmt.Println("[USAGE]: ./TCPChat $port")
			os.Exit(0)
		}
		// testFirstServer()
	} else {
		if *nextPort == "" {
			newPort(&defaultPort)
			port = defaultPort
		} else {
			port = *nextPort
		}
	}
	for _, portUsed := range portsUsed {
		if port == portUsed {
			return
		}
	}
	portsUsed = append(portsUsed, port)
	if !firstTime {
		portWaiting <- true
	}
}

// Increment the default port of 1
func newPort(port *string) {
	var newPortInt int
	var newPort string
	for _, u := range *port {
		newPortInt = newPortInt*10 + int(u-'0')
	}
	newPortInt++
	newPort = fmt.Sprintf("%d",  newPortInt)
	*port = newPort
}

// func testFirstServer() {
// 	timer := time.NewTicker(1*time.Second)
// 	server := *serverNC.NewServer("0.0.0.0:"+port, serverRank)
// 	go func() {
// 		if server.Start() != nil {
// 			fmt.Println("invalid port")
// 			os.Exit(0)
// 		}
// 	}()
// 	<- timer.C
// 	server.Quitch <- true
// 	<- timer.C
// 	timer.Stop()
// }