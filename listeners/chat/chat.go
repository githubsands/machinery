package chat

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var TOKEN = "XXXXXX"

type DiscordMessageHandler func(*discordgo.Session, *discordgo.MessageCreate)

// NewChat creates a new chat given discord message handlers.
// The user of the machinery package is responsible for creating their own discord messages
func NewChat(fs []func(*discordgo.Session, *discordgo.MessageCreate)) *Chat {
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

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

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

	return &Chat{session: discord}
}

type Chat struct {
	session *discordgo.Session
}

func (c *Chat) Run() {
	err := c.session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		panic("")
	}
}
