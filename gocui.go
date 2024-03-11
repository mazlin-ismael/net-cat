package main

import(
	"fmt"
	"strings"
	"github.com/jroimartin/gocui"
)

// loaded all time
func layout(g *gocui.Gui) error {
	// commands for the admin
	maxX, maxY := g.Size()
	v, err := g.SetView("help", maxX-25, 0, maxX-1, 5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "KEYBINDINGS")
		fmt.Fprintln(v, "← →: Move Client")
		fmt.Fprintln(v, "↑ ↓: Move Server")
		fmt.Fprintln(v, "^C: Exit")
	}

	// init view for kick
	vKick, errK := g.SetView("vKick", maxX-25, 6, maxX-1, 8)
	if errK != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vKick.BgColor = gocui.ColorRed
		fmt.Fprintln(vKick, "<- CLICK TO KICK")
	}

	// to permit of the size terminal to be responsive
	_, errTop := g.SetView(serversViews[currentServer].Users[currentView].ViewName, 0, 0, maxX-30, maxY-1)
	if errTop != nil {
		if errTop != gocui.ErrUnknownView {
			return errTop
		}
	}

	// indication to launch a new server with a port
	serverLabel, errServ := g.SetView("serverLabel", maxX-25, 9, maxX-1, 11)
	if errServ != nil {
		if errServ != gocui.ErrUnknownView {
			return errServ
		}
	}
	fmt.Fprintln(serverLabel, "New Server: Port ↓")

	// init the view Editable for write a port to launch a new Server
	vNewServer, errServ := g.SetView("newServer", maxX-25, 12, maxX-1, 14)
	if errServ != nil {
		if errServ != gocui.ErrUnknownView {
			return errServ
		}
	}
	vNewServer.Editable = true
	g.SetCurrentView(vNewServer.Name())
	return nil
}

// init features of the UI
func initGocui(g *gocui.Gui) {
	g.Cursor = true
	g.Mouse = true
	g.SetManagerFunc(layout)
}

// keyBindings
func initKeybindings(g *gocui.Gui) error {
	// Control+C leave the ui and close the servers
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	// ArrowRight | Move on the right view
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nextView(g)
		}); err != nil {
		return err
	}

	// ArrowLeft | Move on the left view
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return preView(g)
		}); err != nil {
		return err
	}

	// ArrowUp | Move on the next Server
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nextServer(g)
		}); err != nil {
		return err
	}

	// ArrowDown | Move on the precedent Server
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return preServer(g)
		}); err != nil {
		return err
	}

	// click on vKick kick the user of the view displayed on the screen
	if err := g.SetKeybinding("vKick", gocui.MouseLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return kickUser(g)
		}); err != nil {
		return err
	}

	// launch a new Server with the entry user
	if err := g.SetKeybinding("newServer", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newServer(g)
		}); err != nil {
		return err
	}
	return nil
}

// launch a new Server with the port given by the user
// if no characters or spaces launch a server with a default port
// if the port is incorrect or already used, nothing is done
func newServer(g *gocui.Gui) error {
	var portTUI = make([]byte, 2048)
	newServer, _ := g.View("newServer")
	n, _ := newServer.Read(portTUI)
	if n == 0 {
		n = 1
	}
	nextPort := strings.TrimSpace(string(portTUI[:n-1]))
	newServer.Clear()
	listenStdout(&nextPort)
	return nil
}

// Move on the next Server
func nextServer(g *gocui.Gui) error{
	if currentServer+1 == len(serversViews) {
		currentServer = 0
	} else {
		currentServer++
	}
	currentView = 0
	g.SetViewOnTop(serversViews[currentServer].Users[currentView].ViewName)
	return nil
}

// Move on the precedent Server
func preServer(g *gocui.Gui) error{
	if currentServer-1 == -1 {
		currentServer = len(serversViews)-1
	} else {
		currentServer--
	}
	currentView = 0
	g.SetViewOnTop(serversViews[currentServer].Users[currentView].ViewName)
	return nil
}

// Move on the next View conn of the server
func nextView(g *gocui.Gui) error {
	if currentView+1 == len(serversViews[currentServer].Users) {
		currentView = 0
	} else {
		currentView++
	}
	g.SetViewOnTop(serversViews[currentServer].Users[currentView].ViewName)
	return nil
}

// Move on the predecent View conn of the server
func preView(g *gocui.Gui) error {
	if currentView-1 == -1 {
		currentView = len(serversViews[currentServer].Users)-1
	} else {
		currentView--
	}
	g.SetViewOnTop(serversViews[currentServer].Users[currentView].ViewName)
	return nil
}

// make a new view
func newView(g *gocui.Gui, viewName string, rank int) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(viewName, 0, 0, maxX-30, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if _, err := g.SetViewOnTop(serversViews[currentServer].Users[currentView].ViewName); err != nil {
		return err
	}
	if strings.HasPrefix(viewName, "server") {
		v.Title = viewName + " port:" + port
	} else {		
		v.Title = "no init name"
	}
	return nil
}

// kick a user
func kickUser(g *gocui.Gui) error {
	if currentView != 0 {
		user := &serversViews[currentServer].Users[currentView]
		server := &serversViews[currentServer].server
		server.TransmettingMSG(server.InitMSG(user.Conn, "kick", nil))
		delete(serversViews[currentServer].server.Clients, user.Conn)
		currentView--
	} 
	return nil
}

// write the logs in the main view of the server
func writeLogs(g *gocui.Gui, server int, log,  direction string) {
	var directionLog string
	v, _ := g.View("server" + fmt.Sprintf("%d", server+1))
	switch direction {
	case "entrance":
		directionLog = "\x1b[32m [ENTER] \x1b[0m"
	case "kick":
		directionLog = " \x1b[41m[KICK]\x1b[0m  "
	case "rename":
		directionLog = " \x1b[36m[RENAME]\x1b[0m "
	case "exit":
		directionLog = " \x1b[31m[EXIT]\x1b[0m "
	}

	v.Write([]byte(directionLog + log))
	g.Update(func(g *gocui.Gui) error {
		if currentView == 0 && currentServer == server {
			g.SetViewOnTop("server" + fmt.Sprintf("%d", server+1))
		}
		return nil
	})
}

// write the logs of the conn in this view
func clientMSG(g *gocui.Gui, server int, user, newLog string) {
	v, _ := g.View(user)
	v.Write([]byte(newLog))
	g.Update(func(g *gocui.Gui) error {
		if  serversViews[currentServer].Users[currentView].ViewName == user && currentServer == server {
			g.SetViewOnTop(user)
		}
		return nil
	})
}