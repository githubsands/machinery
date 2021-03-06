package machinery

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/githubsands/machinery/listeners"
	"github.com/githubsands/machinery/listeners/chat"
	"github.com/githubsands/machinery/listeners/debug"
	"github.com/githubsands/machinery/listeners/grpc"
	"github.com/githubsands/machinery/listeners/observability"
)

// Machinery represents the core machinery of this process. It registers other listeners, controls them, and waits for signals for shutdown
type Machinery struct {
	logger    *log.Logger          // TODO: this should be an interface so user can pick their prefered loggger
	listeners []listeners.Listener // TODO: build this out

	wg       *sync.WaitGroup
	shutdown chan struct{}
	done     chan struct{}

	signals chan os.Signal
	on      bool // represents the state of the machine. if turned off machine is in the process of shutting down

	chatMsgs chan chat.ChatMsg
}

func NewMachinery(cfg Config, l *log.Logger, gss []*grpc.GRPCServer, f []func(*discordgo.Session, *discordgo.MessageCreate)) (*Machinery, chan struct{}) {
	done := make(chan struct{})

	m := &Machinery{wg: &sync.WaitGroup{},
		logger:   l,
		shutdown: make(chan (struct{})),

		on:      true,
		done:    done,
		signals: make(chan os.Signal, 10),

		listeners: make([]listeners.Listener, 0),
	}

	m.addListeners(cfg, gss, f)
	return m, done
}

func (m *Machinery) addListeners(cfg Config, gs []*grpc.GRPCServer, f []func(*discordgo.Session, *discordgo.MessageCreate)) {
	if f != nil {
		chat, chatMsgs := chat.NewChat(f, cfg.ChatHost)
		m.listeners = append(m.listeners, chat)
		m.chatMsgs = chatMsgs
	}

	debug := debug.NewDebug()
	observer := observability.NewObserver()

	if len(gs) > 0 && gs != nil {
		for _, v := range gs {
			m.registerWithChat(v)
			m.listeners = append(m.listeners, v)
		}
	}

	m.listeners = append(m.listeners, debug, observer)
}

func (m *Machinery) SendChatMsg(c chat.ChatMsg) {
	m.chatMsgs <- c
}

func (m *Machinery) Notify() {
	signal.Notify(m.signals, os.Kill, os.Interrupt)
}

func (m *Machinery) registerWithChat(s listeners.ListenerWithChat) {
	if m.chatMsgs != nil {
		lc := listeners.ListenerWithChat(s)
		lc.AddChat(m.chatMsgs)
	}

	return
}

// Run runs all registered listeners
func (m *Machinery) Run() {
	for _, v := range m.listeners {
		m.wg.Add(1)
		go v.Run()
	}
}

// WaitForExit runs in a loop until the conditional is met (1) OS signals are received, (2) the done chan receives a signal
func (m *Machinery) WaitForExit() {
	for m.on = false; !m.on; {
		select {
		case s := <-m.signals:
			m.logger.Printf("Signaled received, %v, exiting the application", s)
			m.on = false
		case <-m.done:
			m.logger.Printf("Done channel received exiting the application")
			m.on = false
		}
	}
}

// WaitForGroup waits for all channels under this wait group
func (m *Machinery) WaitForGroup() {
	m.wg.Wait()
}
