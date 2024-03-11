package server

import (
	"fmt"
	"time"
)

// define the content of the message depending on the typeMGS
func (msg *Message) definePayload(content []byte) {
	switch msg.TypeMSG {
	case "client":
		msg.Payload = msg.clientMSG(content)
	case "entrance":
		msg.Payload = msg.entranceMSG()
	case "exit":
		msg.Payload = msg.exitMSG()
	case "rename":
		msg.Payload = msg.renameMSG()
	case "kick":
		msg.Payload = msg.kickMSG()
	}
}

// for a client MSG
func (msg Message) clientMSG(content []byte) []byte {
	return []byte(fmt.Sprintf("[%s][%s]%s\n", time.Now().Format(time.DateTime), msg.From.Username, content))
}

// for a entrance MSG
func (msg Message) entranceMSG() []byte {
	return append(msg.From.Username, []byte(" has join the chat...\n")...)
}

// for a exit MSG
func (msg Message) exitMSG() []byte {
	return append(msg.From.Username, []byte(" has left the chat...\n")...)
}

// for a rename of a user
func (msg Message) renameMSG() []byte {
	return []byte(fmt.Sprintf("%s has renamed as %s...\n", msg.From.OldName, msg.From.Username))
}

// for a kick of a user
func (msg Message) kickMSG() []byte {
	return append(msg.From.Username, []byte(" has been kicked by and admin...\n")...)
}
