package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	ln, errLn := net.Listen("tcp", ":8989")
	if errLn != nil {
		fmt.Println("errLn", errLn)
		os.Exit(1)
	}
	conn, errConn := ln.Accept()
	if errConn != nil {
		fmt.Println("errConn", errConn)
	}
	HandlerConnection(conn)
	wg.Wait()
}

func HandlerConnection(conn net.Conn) {
	stdout, _ := io.Copy(os.Stdout, conn)
	fmt.Println(stdout)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("err", err)
			conn.Close()
			os.Exit(1)
		}
		fmt.Printf("Received %s\n", buffer[:n])
	}

}