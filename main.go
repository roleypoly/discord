package main // import "github.com/roleypoly/discord"

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}

	discord.AddHandler(messageHandler)

	err = discord.Open()
	if err != nil {
		log.Fatalln(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func findInUsers(ul []*discordgo.User, pred func(*discordgo.User) bool) *discordgo.User {
	for _, u := range ul {
		if pred(u) {
			return u
		}
	}

	return nil
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	if findInUsers(m.Mentions, func(u *discordgo.User) bool {
		return u.ID == s.State.User.ID
	}) == nil {
		return
	}

	s.ChannelMessageSend(m.ChannelID, `:beginner: Assign your roles here! https://roleypoly.com/s/`+m.GuildID)
}
