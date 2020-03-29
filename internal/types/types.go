package types

import (
	discordgobot "github.com/lampjaw/discordclient"
	"regexp"
)

type Message struct {
	Message discordgobot.Message
	GuildID string
	RawText string
}

type CommandHandler func(message Message) string

type Command struct {
	Matcher  *regexp.Regexp
	Callback CommandHandler
}
