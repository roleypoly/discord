package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/roleypoly/discord/responders/dm"
	"github.com/roleypoly/discord/responders/text"
)

func idInList(list []string, id string) bool {
	for _, v := range list {
		if v == id {
			return true
		}
	}

	return false
}

func findInUsers(ul []*discordgo.User, pred func(*discordgo.User) bool) *discordgo.User {
	for _, u := range ul {
		if pred(u) {
			return u
		}
	}

	return nil
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	if m.GuildID == "" {
		handleDM(s, m)
	} else {
		handleText(s, m)
	}
}

func handleDM(s *discordgo.Session, m *discordgo.MessageCreate) {
	dm.DMResponder(s, m)
}

func handleText(s *discordgo.Session, m *discordgo.MessageCreate) {
	if findInUsers(m.Mentions, func(u *discordgo.User) bool {
		return u.ID == s.State.User.ID
	}) == nil {
		return
	}

	if !idInList(rootUsers, m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, `:beginner: Assign your roles here! https://roleypoly.com/s/`+m.GuildID)
	} else {
		text.TextResponder(s, m)
	}
}
