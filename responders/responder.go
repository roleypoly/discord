package responders

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Match       *regexp.Regexp
	Description string
	Handler     func(matches [][]string, s *discordgo.Session, m *discordgo.MessageCreate) string
}

func (c Command) Test(msg string) bool {
	return c.Match.MatchString(msg)
}

func (c Command) Run(msg string, s *discordgo.Session, m *discordgo.MessageCreate) {
	rv := c.Handler(c.Match.FindAllStringSubmatch(msg, -1), s, m)
	s.ChannelMessageSend(m.ChannelID, rv)
}

type Responder struct {
	Commands []Command
}

func (r *Responder) FindAndExecute(msg string, s *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, cmd := range r.Commands {
		if cmd.Test(msg) {
			cmd.Run(msg, s, m)
			return true
		}
	}

	return false
}
