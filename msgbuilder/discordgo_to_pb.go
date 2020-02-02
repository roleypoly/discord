package msgbuilder

import (
	"github.com/bwmarrin/discordgo"
	pbDiscord "github.com/roleypoly/rpc/discord"
	pbShared "github.com/roleypoly/rpc/shared"
)

func Guild(guild *discordgo.Guild) *pbShared.Guild {
	if guild == nil {
		return nil
	}

	return &pbShared.Guild{
		ID:          guild.ID,
		Name:        guild.Name,
		Icon:        guild.Icon,
		OwnerID:     guild.OwnerID,
		MemberCount: int32(guild.MemberCount),
		Splash:      guild.Splash,
	}
}

func User(user *discordgo.User) *pbShared.DiscordUser {
	if user == nil {
		return nil
	}

	return &pbShared.DiscordUser{
		ID:            user.ID,
		Username:      user.Username,
		Discriminator: user.Discriminator,
		Avatar:        user.Avatar,
		Bot:           user.Bot,
	}
}

func Member(member *discordgo.Member) *pbDiscord.Member {
	if member == nil {
		return nil
	}

	return &pbDiscord.Member{
		GuildID: member.GuildID,
		Roles:   member.Roles,
		Nick:    member.Nick,
		User:    User(member.User),
	}
}

func Roles(roles []*discordgo.Role) []*pbShared.Role {
	protoRoles := make([]*pbShared.Role, len(roles))

	for idx, role := range roles {
		protoRoles[idx] = Role(role)
	}

	return protoRoles
}

func Role(role *discordgo.Role) *pbShared.Role {
	if role == nil {
		return nil
	}

	return &pbShared.Role{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: int64(role.Permissions),
		Color:       int32(role.Color),
		Position:    int32(role.Position),
		Managed:     role.Managed,
	}
}
