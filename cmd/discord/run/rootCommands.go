package run

import (
	"regexp"

	"github.com/lampjaw/discordclient"
)

type Command struct {
	Match *regexp.Regexp
	Fn    func(l *Listener, message discordclient.Message) string
}

var rootCommands []Command = []Command{
	{
		Match: regexp.MustCompile("stats$"),
		Fn:    rootStats,
	},
}

func (l *Listener) handleRoot(message discordclient.Message) bool {
	for _, cmd := range rootCommands {
		if cmd.Match.MatchString(message.RawMessage()) {
			output := cmd.Fn(l, message)

			if !l.readOnly {
				l.Bot.SendMessage(message.Channel(), output)
			}

			return true
		}
	}

	return false
}

func rootStats(l *Listener, message discordclient.Message) string {
	return "stats"
}
