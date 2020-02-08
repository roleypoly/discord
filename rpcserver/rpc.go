package rpcserver

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/protobuf/ptypes/empty"
	discordgobot "github.com/lampjaw/discordclient"
	"github.com/roleypoly/discord/internal/permissions"
	"github.com/roleypoly/discord/msgbuilder"
	pbDiscord "github.com/roleypoly/rpc/discord"
	pbShared "github.com/roleypoly/rpc/shared"
	"k8s.io/klog"
)

// DiscordService is a gRPC implementation of rpc/discord.proto#Discord.
type DiscordService struct {
	pbDiscord.DiscordServer
	Discord *discordgobot.DiscordClient
}

// ListGuilds lists every guild in state.
func (d *DiscordService) ListGuilds(ctx context.Context, req *empty.Empty) (*pbShared.GuildList, error) {
	guildlist := &pbShared.GuildList{}
	botguilds := d.Discord.Guilds()
	guildlist.Guilds = make([]*pbShared.Guild, len(botguilds))

	for idx, guild := range botguilds {
		guildlist.Guilds[idx] = msgbuilder.Guild(guild)
	}

	return guildlist, nil
}

// GetGuild fetches a single Guild from state.
func (d *DiscordService) GetGuild(ctx context.Context, req *pbShared.IDQuery) (*pbShared.Guild, error) {
	g, err := d.Discord.Guild(req.GuildID)
	return msgbuilder.Guild(g), err
}

// GetGuildsByMember searches for guilds that include a certain member.
func (d *DiscordService) GetGuildsByMember(ctx context.Context, req *pbShared.IDQuery) (*pbShared.GuildList, error) {
	memberGuilds := &pbShared.GuildList{}
	for _, guild := range d.Discord.Guilds() {
		mem, err := d.Discord.GuildMember(req.MemberID, guild.ID)
		if err != nil {
			continue
		}

		if mem != nil {
			memberGuilds.Guilds = append(memberGuilds.Guilds, msgbuilder.Guild(guild))
		}
	}

	return memberGuilds, nil
}

// GetMember fetches a guild member by a server.
func (d *DiscordService) GetMember(ctx context.Context, req *pbShared.IDQuery) (*pbDiscord.Member, error) {
	member, err := d.Discord.GuildMember(req.MemberID, req.GuildID)
	return msgbuilder.Member(member), err
}

// UpdateMember
func (d *DiscordService) UpdateMember(ctx context.Context, req *pbDiscord.Member) (*pbDiscord.Member, error) {
	err := d.Discord.Session.GuildMemberEdit(req.GuildID, req.User.ID, req.Roles)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func getRoleFromRoles(roleID string, roles []*discordgo.Role) *discordgo.Role {
	for _, role := range roles {
		if role.ID == roleID {
			return role
		}
	}

	return nil
}

func (d *DiscordService) calculateSafety(guild *discordgo.Guild, role *pbShared.Role) pbShared.Role_RoleSafety {
	if permissions.ProtoRoleHasPermission(role, discordgo.PermissionManageRoles) || permissions.ProtoRoleHasPermission(role, discordgo.PermissionAdministrator) {
		return pbShared.Role_dangerousPermissions
	}

	ownMember, err := d.Discord.GuildMember(d.Discord.UserID(), guild.ID)
	if err != nil {
		klog.Error("ownMember is nil in", guild.ID)
		return pbShared.Role_higherThanBot
	}

	highestOwnRolePosition := 0
	for _, roleID := range ownMember.Roles {
		checkRole := getRoleFromRoles(roleID, guild.Roles)

		if checkRole.Position > highestOwnRolePosition {
			highestOwnRolePosition = checkRole.Position
		}
	}

	if role.Position > int32(highestOwnRolePosition) {
		return pbShared.Role_higherThanBot
	}

	return pbShared.Role_safe
}

func (d *DiscordService) GetGuildRoles(ctx context.Context, req *pbShared.IDQuery) (*pbShared.GuildRoles, error) {
	guild, err := d.Discord.Guild(req.GuildID)
	if err != nil {
		return nil, err
	}

	roles := msgbuilder.Roles(guild.Roles)

	for _, role := range roles {
		role.Safety = d.calculateSafety(guild, role)
	}

	return &pbShared.GuildRoles{
		ID:    guild.ID,
		Roles: roles,
	}, nil
}

func (d *DiscordService) GetUser(ctx context.Context, req *pbShared.IDQuery) (*pbShared.DiscordUser, error) {
	user, err := d.Discord.Session.User(req.MemberID)
	return msgbuilder.User(user), err
}

func (d *DiscordService) OwnUser(ctx context.Context, req *empty.Empty) (*pbShared.DiscordUser, error) {
	return d.GetUser(ctx, &pbShared.IDQuery{MemberID: d.Discord.UserID()})
}
