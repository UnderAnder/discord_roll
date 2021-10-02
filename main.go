package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	rand.Seed(time.Now().Unix())

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!roll" {

		s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" :game_die: "+strconv.Itoa(rand.Intn(101)))
	}

	if m.Content == "!bottle" {
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			fmt.Println("error can't get guild members", err)
			return
		}

		randMember := members[rand.Intn(len(members))]
		for randMember.User.ID == m.Author.ID {
			randMember = members[rand.Intn(len(members))]
		}

		s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" :kiss: "+getMemberNick(randMember))
	}
}

func getMessageAuthorNick(m *discordgo.MessageCreate) string {
	if m.Member.Nick != "" {
		return m.Member.Nick
	} else {
		return m.Author.Username
	}
}

func getMemberNick(m *discordgo.Member) string {
	if m.Nick != "" {
		return m.Nick
	} else {
		return m.User.Username
	}
}
