package rpcserver

import "github.com/bwmarrin/discordgo"

var testMember = &discordgo.Member{
	Roles: []string{
		"bot",
		"color-blue",
	},
}

var testGuild = &discordgo.Guild{
	Roles: []*discordgo.Role{
		{
			ID:          "admin",
			Permissions: discordgo.PermissionAdministrator,
			Position:    10,
		},
		{
			ID:          "unprivileged-and-higher",
			Permissions: 0,
			Position:    7,
		},
		{
			ID:          "bot",
			Permissions: discordgo.PermissionManageRoles | discordgo.PermissionSendMessages,
			Position:    5,
		},
		{
			ID:          "mod",
			Permissions: discordgo.PermissionKickMembers | discordgo.PermissionManageRoles | discordgo.PermissionManageMessages,
			Position:    3,
		},
		{
			ID:          "color-red",
			Permissions: 0,
			Position:    2,
		},
		{
			ID:          "color-blue",
			Permissions: 0,
			Position:    1,
		},
	},
}
