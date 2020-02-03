package run

import (
	"log"
	"os"
	"regexp"

	discordgobot "github.com/lampjaw/discordclient"
	"github.com/roleypoly/discord/internal/strings"
)

// Listener is a message processor built on top of the discord client.
type Listener struct {
	Bot         *discordgobot.DiscordClient
	selfMention *regexp.Regexp
	readOnly    bool
}

// Run begins the message processor flow and loop.
func (l *Listener) Run() {
	l.readOnly = os.Getenv("READ_ONLY") == "1"

	msgChan, err := l.Bot.Listen(-1)
	if err != nil {
		log.Fatalln("err: Listener.Run, discord client listen --", err)
		return
	}

	log.Println("shards:", len(l.Bot.Sessions))

	l.selfMention = regexp.MustCompile("<@!?" + l.Bot.UserID() + ">")
	go l.startListening(msgChan)
}

func (l *Listener) startListening(msgChan <-chan discordgobot.Message) {
	log.Println("discord bot running")
	for {
		message := <-msgChan
		go l.handleMessage(message)
	}
}

func (l *Listener) handleMessage(message discordgobot.Message) {
	raw := message.RawMessage()
	guildID, err := message.ResolveGuildID()
	if err != nil {
		log.Println("err: Listener.handleMessage, guildID resolve --", err)
		return
	}

	if guildID == "" {
		// this is a DM
	} else {
		if l.selfMention.MatchString(raw) {
			log.Println("raw", raw)
			log.Printf("guildID -- `%s`\n", guildID)
			response := strings.Render(strings.MentionResponse, strings.MentionResponseData{GuildID: guildID, AppURL: os.Getenv("APP_URL")})

			if !l.readOnly {
				l.Bot.SendMessage(message.Channel(), response)
			}
		}
	}
}
