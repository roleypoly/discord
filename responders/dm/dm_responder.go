package dm

import (
	"bytes"
	"context"
	"log"
	"regexp"
	"strings"
	"text/template"

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

var (
	loginTemplate = template.Must(template.New("logintemplate").Parse(`**Hey there {{.User}}!** <a:promareFlame:624850108667789333>

Use this secret code: **{{.MagicWords}}**
Or click here: <{{.MagicURL}}>

This code will self-destruct in 1 hour.`))
)

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

	tmplData := struct {
		RandomMessage string
		User          string
		MagicWords    string
		MagicURL      string
	}{
		RandomMessage: "Hi~",
		User:          m.Author.Username,
		MagicWords:    authChallenge.MagicWords,
		MagicURL:      authChallenge.MagicUrl,
	}

	buf := bytes.Buffer{}
	err = loginTemplate.Execute(&buf, tmplData)
	if err != nil {
		log.Println(err)
		return ""
	}

	return buf.String()
}
