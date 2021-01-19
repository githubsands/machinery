package machinery

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/githubsands/investmentmanager/machinery/chat"
	"github.com/githubsands/investmentmanager/machinery/observability"
)

// listeners defines what can be registered with machinery
// TODO: make this abide to as many different listeners as possible i.e grpc or http servers
type listener interface {
	Run()
}

// Machinery represents the core machinery of this process. It registers other listeners, controls them, and waits for signals for shutdown
type Machinery struct {
	logger    *log.Logger // TODO: this should be an interface so user can pick their prefered loggger
	listeners []listener  // TODO: build this out

	wg       *sync.WaitGroup
	shutdown chan struct{}
	done     chan struct{}

	signals chan os.Signal
	on      bool // represents the state of the machine. if turned off machine is in the process of shutting down

	// observer observer
}

func NewMachinery(l *log.Logger) (*Machinery, chan struct{}) {
	done := make(chan struct{})

	m := &Machinery{wg: &sync.WaitGroup{},
		logger:   l,
		shutdown: make(chan (struct{})),

		done:      done,
		signals:   make(chan os.Signal, 10),
		listeners: make([]listener, 0),
	}

	m.addListeners()
	return m, done
}

func (m *Machinery) addListeners() {
	observer := observability.NewObserver()
	chat := chat.NewChat()

	m.listeners = append(m.listeners, observer, chat)
}

func (m *Machinery) Notify() {
	signal.Notify(m.signals, os.Kill, os.Interrupt)
}

func (m *Machinery) Register(s listener) {
	m.listeners = append(m.listeners, s)
}

// Run runs all registered listeners
func (m *Machinery) Run() {
	for _, v := range m.listeners {
		m.wg.Add(1)
		go v.Run()
	}
}

func (m *Machinery) off() {
	m.on = true
}

// WaitForExit runs in a loop until the conditional is met (1) OS signals are received, (2) the done chan receives a signal
func (m *Machinery) WaitForExit() {
	for exit := false; !exit; {
		select {
		case s := <-m.signals:
			m.logger.Printf("Signaled received, %v, exiting the application", s)
			m.off()
		case <-m.done:
			m.logger.Printf("Done channel received exiting the application")
			m.off()
		}
	}
}

// WaitForGroup waits for all channels under this wait group
func (m *Machinery) WaitForGroup() {
	m.wg.Wait()
}
