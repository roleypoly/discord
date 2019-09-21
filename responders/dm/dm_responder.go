package dm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/roleypoly/discord/auth"
	"github.com/roleypoly/discord/responders"
	pbAuth "github.com/roleypoly/rpc/auth/backend"
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
	authConn, err := auth.NewAuthConnector()
	if err != nil {
		log.Println(err)
		return ""
	}

	authChallenge, err := authConn.Client.GetSessionChallenge(context.Background(), &pbAuth.UserSlug{
		UserID: m.Author.ID,
	})
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Hey there <@%s>! Tell me `%s` or click %s!", authChallenge.UserID, authChallenge.MagicWords, authChallenge.MagicUrl)
}
