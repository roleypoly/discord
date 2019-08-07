package dm

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/roleypoly/discord/responders"
)

var dmResponder = &responders.Responder{
	Commands: []responders.Command{{
		Match:   regexp.MustCompile(``),
		Handler: handleStats,
	}},
}

func DMResponder(s *discordgo.Session, m *discordgo.MessageCreate) {
	myID := s.State.User.ID
	msg := strings.Replace(m.Content, fmt.Sprintf(`<@$%s>`, myID), "", 1)
	msg = strings.Replace(m.Content, fmt.Sprintf(`<@!$%s>`, myID), "", 1)

	dmResponder.FindAndExecute(msg, s, m)
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
