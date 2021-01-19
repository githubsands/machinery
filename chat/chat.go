package chat

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var TOKEN = "XXXXXX"

func NewChat() *Chat {
	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		panic("")
		fmt.Println("error creating Discord session,", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
