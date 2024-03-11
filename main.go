package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jroimartin/gocui"

	serverNC "NETCAT/SERVER"
)

// struct use for the UI interface
var serversViews []ViewsSettings

func main() {
	// associate the newPort to 8989 by default if go run . | else by the port given by the user
	listenStdout(&port)
	// init the UI interfacce
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	initGocui(g)

	go func() {
		for {
			// append a ViewsSettings to the struct []ViewsSettings
			appendServerGUI()
			// append a view to the UI General for logs of the server
			newView(g, serversViews[serverRank].Users[0].ViewName, serverRank)
			// launch a newServer with port given
			go serverInit("0.0.0.0:"+port, g)
			if firstTime {
				firstTime = false
				initBinds <- true
			}
			// Wait for the next server to launch
			<-portWaiting
			serverRank++
		}
	}()
	// Manage keyBindings of UI and actions
	<-initBinds
	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// New Server
func serverInit(ListenAddr string, g *gocui.Gui) {
	// Init new Server
	server := *serverNC.NewServer(ListenAddr, serverRank)
	serversViews[server.Rank].server = server

	go func() {
		// logs for the server
		file, _ := os.Create("LOGS/" + server.ListenAddr + ".log")
		for msg := range server.Msgch {
			// write logs
			file.Write(msg.Payload)

			if msg.TypeMSG != "client" {
				// write logs in the general view of the server
				writeLogs(g, server.Rank, string(msg.Payload), msg.TypeMSG)
			}

			switch msg.TypeMSG {
			case "exit", "kick":
				// delete the user from the usersView and delete view
				for i, u := range serversViews[server.Rank].Users {
					if u.ViewName == msg.From.ViewName {
						if i == len(serversViews[server.Rank].Users)-1 {
							serversViews[server.Rank].Users = serversViews[server.Rank].Users[:i]
						} else {
							serversViews[server.Rank].Users = append(serversViews[server.Rank].Users[:i], serversViews[server.Rank].Users[i+1:]...)
						}
						if currentView >= i && server.Rank == currentServer {
							currentView--
						}
						g.DeleteView(u.ViewName)
						break
					}
				}
			case "newView":
				// Make a new view for the user
				serversViews[server.Rank].Users = append(serversViews[server.Rank].Users, struct {
					Conn     net.Conn
					Name     string
					ViewName string
				}{
					Name:     string(msg.From.Username),
					Conn:     msg.From.Conn,
					ViewName: msg.From.ViewName,
				})
				newView(g, msg.From.ViewName, server.Rank)
				view, _ := g.View(msg.From.ViewName)
				view.Title = msg.From.ViewName

			case "entrance", "rename":
				// change the title of the view
				for _, client := range serversViews[server.Rank].Users {
					if client.Conn == msg.From.Conn {
						v, _ := g.View(msg.From.ViewName)
						v.Title = string(msg.From.Username)
					}
				}
			}

			for conn, user := range server.Clients {
				// write msg from a user in the others users terminal, and write logs from user in his view
				if conn != msg.From.Conn {
					if user.Username == nil {
						continue
					}
					conn.Write([]byte(msg.Payload))
				} else if msg.TypeMSG == "client" && msg.TypeMSG != "kick" {
					conn.Write([]byte("\033[1A\033[K" + string(msg.Payload)))
					clientMSG(g, server.Rank, msg.From.ViewName, string(msg.Payload))
				}
			}
		}
	}()

	// Launch new Server
	err := server.Start()
	if err != nil {
		// if err delete server from the views
		g.DeleteView("server" + fmt.Sprintf("%d", serverRank+1))
		serversViews = serversViews[:serverRank]
		serverRank--
		return
	}
}

// append a ViewsSettings to the struct []ViewsSettings with server+rank
func appendServerGUI() {
	serversViews = append(serversViews, ViewsSettings{})
	serversViews[serverRank].Users = append(serversViews[serverRank].Users, struct {
		Conn     net.Conn
		Name     string
		ViewName string
	}{
		ViewName: "server" + fmt.Sprintf("%d", serverRank+1),
	})
}
