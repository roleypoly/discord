package text

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/roleypoly/discord/responders"
)

var textResponder = &responders.Responder{
	Commands: []responders.Command{{
		Match:   regexp.MustCompile(`stats`),
		Handler: handleStats,
	}, {
		Match: regexp.MustCompile(`say (.*)$`),
		Handler: func(matches [][]string, s *discordgo.Session, m *discordgo.MessageCreate) string {
			return matches[0][1]
		},
	}},
}

func TextResponder(s *discordgo.Session, m *discordgo.MessageCreate) {
	myID := s.State.User.ID
	msg := strings.Replace(m.Content, fmt.Sprintf(`<@%s>`, myID), "", 1)
	msg = strings.Replace(m.Content, fmt.Sprintf(`<@!%s>`, myID), "", 1)

	if !textResponder.FindAndExecute(msg, s, m) {
		SendDefaultResponse(s, m)
	}
}

func SendDefaultResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, `:beginner: Assign your roles here! `+os.Getenv("APP_URL")+`/s/`+m.GuildID)
}

func memberCountFromGuilds(gs []*discordgo.Guild) (count int) {
	for _, g := range gs {
		count += g.MemberCount
	}

	return
}

func roleCountFromGuilds(gs []*discordgo.Guild) (count int) {
	for _, g := range gs {
		count += len(g.Roles)
	}

	return
}

func handleStats(_ [][]string, s *discordgo.Session, m *discordgo.MessageCreate) string {
	guilds := len(s.State.Guilds)
	users := memberCountFromGuilds(s.State.Guilds)
	roles := roleCountFromGuilds(s.State.Guilds)

	return fmt.Sprintf(`**Stats** :chart_with_upwards_trend:

:couple_ww: **Users Served:** %d
:beginner: **Servers:** %d
:white_flower: **Roles Seen:** %d`, users, guilds, roles)

}