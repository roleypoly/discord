package run

import (
	"regexp"

	"github.com/lampjaw/discordclient"
	"github.com/roleypoly/discord/internal/strings"
	"github.com/roleypoly/discord/internal/version"
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
	stats := strings.RootStatsData{}

	// People Info
	stats.Users = l.Bot.UserCount()
	stats.Guilds = len(l.Bot.Guilds())
	roles := 0
	for _, guild := range l.Bot.Guilds() {
		roles += len(guild.Roles)
	}
	stats.Roles = roles

	// Bot Info
	stats.Shards = len(l.Bot.Sessions)
	stats.BuildDate = version.BuildDate
	stats.GitBranch = version.GitBranch
	stats.GitCommit = version.GitCommit

	return strings.Render(strings.RootStats, stats)
}
