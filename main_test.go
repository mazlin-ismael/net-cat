package main

import (
	serverNC "NETCAT/SERVER"
	"net"
	"testing"
	"time"
)

var wait chan bool = make(chan bool)
var end = time.After(4*time.Second)

// server init and start and conn connected
func TestMain(t *testing.T) {
	server := serverNC.NewServer("localhost:1234", 1)
	go server.Start()
	ta := time.NewTicker(100*time.Millisecond)
	<- ta.C
	for range make([]struct{}, 11) {
		go func() {
			_, err := net.Dial("tcp", "localhost:1234")
			if err != nil {
				t.Fail()
			}
		}()
	}
	go func() {
		<- end
		wait <- true
	}()
	<- wait
	if len(server.Clients) > 10 {
		t.Fail()
	}
}
