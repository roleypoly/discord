package dm

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/roleypoly/discord/responders"
)

var dmResponder = &responders.Responder{
	Commands: []responders.Command{{
		Match:   regexp.MustCompile(`((log|sign) ?in|auth)`),
		Handler: handleAuth,
	}},
}

func DMResponder(s *discordgo.Session, m *discordgo.MessageCreate) {
	dmResponder.FindAndExecute(strings.ToLower(m.Content), s, m)
}

func handleAuth(_ [][]string, s *discordgo.Session, m *discordgo.MessageCreate) string {
	return `This function is currently disabled.`
}
