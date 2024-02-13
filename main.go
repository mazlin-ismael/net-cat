package main

import (
	"app/server"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

func main() {
	var port string = "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		return
	}

	// ip 192.168.100.224 campus tech
	var listenAddr string = "0.0.0.0:"+port
	server := serverNC.NewServer(listenAddr)
	fmt.Println("Server launching on: " + listenAddr)


	go func() {
		file, _ := os.Create("logs.txt")
		for  msg := range server.Msgch {
			var newLog string
			newLog += "["+time.Now().Format(time.DateTime)+"][" + string(msg.From[:len(msg.From)-1]) + "]" + string(msg.Payload)
			file.WriteString(newLog)
			for _, client := range server.Clients {
				if !reflect.DeepEqual(msg.From, client.UserName) {
					client.Conn.Write([]byte("["+time.Now().Format(time.DateTime)+"]["))
					client.Conn.Write(msg.From[:len(msg.From)-1])
					client.Conn.Write([]byte("]:"))
					client.Conn.Write(msg.Payload)
					
				}
			}
			// fmt.Printf("received message from connection (%s): %s\n", msg.From, string(msg.Payload))
		}
	}()

	log.Fatal(server.Start())
}