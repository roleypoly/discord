package run

import (
	"fmt"
	"regexp"

	"github.com/lampjaw/discordclient"
	stringrenderer "github.com/roleypoly/discord/internal/strings"
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
	{
		Match: regexp.MustCompile("shard of [0-9]+$"),
		Fn:    rootGetShard,
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
	stats := stringrenderer.RootStatsData{}

	memberCount := 0
	for _, guild := range l.Bot.Guilds() {
		memberCount += guild.MemberCount
	}

	// People Info
	stats.Users = memberCount
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

	return stringrenderer.Render(stringrenderer.RootStats, stats)
}

var shardMatch = regexp.MustCompile("shard of ([0-9]+)$")

func rootGetShard(l *Listener, message discordclient.Message) string {
	id := shardMatch.FindAllStringSubmatch(message.RawMessage(), 1)[0][1]

	if l.Bot.Session.ShardCount == 1 {
		session := l.Bot.Session
		guild, err := session.State.Guild(id)
		if guild != nil || err == nil {
			return fmt.Sprintf("Shard of **%s** is **%d** (of %d)", guild.Name, session.ShardID+1, session.ShardCount)
		}
	}

	for _, session := range l.Bot.Sessions {
		guild, err := session.State.Guild(id)
		if guild != nil || err == nil {
			return fmt.Sprintf("Shard of **%s** is **%d** (of %d)", guild.Name, session.ShardID+1, session.ShardCount)
		}
	}

	return "Shard not found."
}
