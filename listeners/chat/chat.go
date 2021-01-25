package chat

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var TOKEN = "XXXXXX"

type DiscordMessageHandler func(*discordgo.Session, *discordgo.MessageCreate)

// NewChat creates a new chat given discord message handlers.
// The user of the machinery package is responsible for creating their own discord messages
func NewChat(fs []func(*discordgo.Session, *discordgo.MessageCreate)) (*Chat, chan ChatMsg) {
	if len(fs) == 0 {
		panic("")
	}

	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		panic("")
		fmt.Println("error creating Discord session,", err)
	}
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		panic("")
	}

	for _, v := range fs {
		dg.AddHandler(v)
	}

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		panic("")
	}

	discord, err := discordgo.New("Bot " + "authentication token")
	if err != nil {
		panic("")
	}

	chat := new(Chat)
	chat.session = discord
	chat.msgs = make(chan ChatMsg, 1)

	// start receiving messages
	go chat.receiveMsgs()

	return chat, chat.msgs
}

type Chat struct {
	msgs    chan ChatMsg
	session *discordgo.Session
}

func (c *Chat) Run() {
	err := c.session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		panic("")
	}
}

// receiveMsgs gets msgs from all listeners and routes them to their respective update
func (c *Chat) receiveMsgs() {
	for {
		select {
		case msg := <-c.msgs:
			switch v := msg.(type) {
			case StatusMsg:
				c.updateStatus(v)
			case StreamingStatusMsg:
				c.updateStreamingStatus(v)
			case ListeningStatusMsg:
				c.updateListeningStatus(v)
			case StatusComplexMsg:
				c.updateStatusComplex(v)
			default:
				fmt.Printf("Undefined msg type: %v", msg)
			}
		default:
			fmt.Printf("Error getting message: %v")
		}
	}
}

func (c *Chat) updateStatus(msg StatusMsg) error {
	return c.session.UpdateStatus(msg.i, msg.msg)
}

// TODO: update this to abide by real info streaming status needs
func (c *Chat) updateStreamingStatus(msg StreamingStatusMsg) error {
	return c.session.UpdateStreamingStatus(msg.i, msg.msg, msg.msg)
}

func (c *Chat) updateListeningStatus(msg ListeningStatusMsg) error {
	return c.session.UpdateListeningStatus(msg.msg)
}

func (c *Chat) updateStatusComplex(msg StatusComplexMsg) error {
	d := discordgo.UpdateStatusData{}
	return c.session.UpdateStatusComplex(d)
}
