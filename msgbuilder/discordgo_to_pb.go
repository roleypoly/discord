package msgbuilder

import (
	"github.com/bwmarrin/discordgo"
	pb "github.com/roleypoly/rpc/discord"
)

func Guild(guild *discordgo.Guild) *pb.Guild {
	if guild == nil {
		return nil
	}

	return &pb.Guild{
		ID:          guild.ID,
		Name:        guild.Name,
		Icon:        guild.Icon,
		OwnerID:     guild.OwnerID,
		MemberCount: int32(guild.MemberCount),
		Splash:      guild.Splash,
	}
}

func User(user *discordgo.User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		ID:            user.ID,
		Username:      user.Username,
		Discriminator: user.Discriminator,
		Avatar:        user.Avatar,
		Bot:           user.Bot,
	}
}

func Member(member *discordgo.Member) *pb.Member {
	if member == nil {
		return nil
	}

	return &pb.Member{
		GuildID: member.GuildID,
		Roles:   member.Roles,
		Nick:    member.Nick,
		User:    User(member.User),
	}
}

func Roles(roles []*discordgo.Role) []*pb.Role {
	protoRoles := make([]*pb.Role, len(roles))

	for idx, role := range roles {
		protoRoles[idx] = Role(role)
	}

	return protoRoles
}

func Role(role *discordgo.Role) *pb.Role {
	if role == nil {
		return nil
	}

	return &pb.Role{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: int64(role.Permissions),
		Color:       int32(role.Color),
		Position:    int32(role.Position),
		Managed:     role.Managed,
	}
}
