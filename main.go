package main

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var allowedGuild []string = []string{"683668323095019522", "843780677019500565"} // Private guilds.

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages
	dg.Identify.Presence.Status = string(discordgo.StatusOffline)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening webdocket connection: ", err)
		return
	}

	err = dg.UpdateCustomStatus("Message moderation.")
	if err != nil {
		fmt.Println("Error setting bot status: ", err)
		return
	}

	dg.Identify.Presence.Status = string(discordgo.StatusDoNotDisturb)

	// Fetch information about the application.
	user, err := dg.User("@me")
	if err != nil {
		fmt.Println("Error fetching bot user information: ", err)
		return
	}

	fmt.Println(user.Username + " is now connected.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore messages came from non-allowed guilds.
	if !slices.Contains(allowedGuild, m.GuildID) {
		return
	}

	// Check if message contains toxicity.
	// If it is, check if probability is more than 70%.
	tox, err := Toxicity(m.Content)
	if err != nil {
		fmt.Println("Error using ML algorithms: ", err)
		return
	}

	if tox.Toxic > tox.Neutral && tox.Toxic > 0.7 {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}
