package listeners

import "github.com/githubsands/machinery/listeners/chat"

// listeners defines what can be registered with machinery
// TODO: make this abide to as many different listeners as possible i.e grpc or http servers

type ListenerWithChat interface {
	AddChat(chan<- chat.ChatMsg)
	SendToChat(chat.ChatMsg)

	Listener
}

type Listener interface {
	Run()
}
