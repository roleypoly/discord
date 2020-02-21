package run

import (
	"os"
	"regexp"

	discordgobot "github.com/lampjaw/discordclient"
	"github.com/roleypoly/discord/internal/strings"
	"k8s.io/klog"
)

// Listener is a message processor built on top of the discord client.
type Listener struct {
	Bot          *discordgobot.DiscordClient
	selfMention  *regexp.Regexp
	readOnly     bool
	RootUsers    []string
	BotWhitelist []string
}

// Run begins the message processor flow and loop.
func (l *Listener) Run() {
	l.readOnly = os.Getenv("READ_ONLY") == "1"

	msgChan, err := l.Bot.Listen(-1)
	if err != nil {
		klog.Exit("err: Listener.Run, discord client listen --", err)
		return
	}

	klog.Info("shards:", len(l.Bot.Sessions))

	l.selfMention = regexp.MustCompile("<@!?" + l.Bot.UserID() + ">")
	go l.startListening(msgChan)
}

func (l *Listener) startListening(msgChan <-chan discordgobot.Message) {
	klog.Info("discord bot running")
	for {
		message := <-msgChan
		go l.handleMessage(message)
	}
}

func (l *Listener) isRoot(userID string) bool {
	for _, id := range l.RootUsers {
		if userID == id {
			return true
		}
	}

	return false
}

func (l *Listener) isWhitelistedBot(userID string) bool {
	for _, id := range l.BotWhitelist {
		if userID == id {
			return true
		}
	}

	return false
}

func (l *Listener) handleMessage(message discordgobot.Message) {
	if message.IsBot() && !l.isWhitelistedBot(message.UserID()) {
		return
	}

	raw := message.RawMessage()
	guildID, err := message.ResolveGuildID()
	if err != nil {
		klog.Exit("err: Listener.handleMessage, guildID resolve --", err)
		return
	}

	if guildID == "" {
		// this is a DM
	} else {
		if l.selfMention.MatchString(raw) {
			if l.isRoot(message.UserID()) {
				ok := l.handleRoot(message)
				if ok {
					return
				}
			}

			response := strings.Render(strings.MentionResponse, strings.MentionResponseData{GuildID: guildID, AppURL: os.Getenv("APP_URL")})

			if !l.readOnly {
				l.Bot.SendMessage(message.Channel(), response)
			}
		}
	}
}
