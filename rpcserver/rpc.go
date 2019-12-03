package rpcserver

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	discordgobot "github.com/lampjaw/discordclient"
	"github.com/roleypoly/discord/msgbuilder"
	pb "github.com/roleypoly/rpc/discord"
)

// DiscordService is a gRPC implementation of rpc/discord.proto#Discord.
type DiscordService struct {
	pb.DiscordServer
	Discord *discordgobot.DiscordClient
}

// ListGuilds lists every guild in state.
func (d *DiscordService) ListGuilds(ctx context.Context, req *empty.Empty) (*pb.GuildList, error) {
	guildlist := &pb.GuildList{}
	botguilds := d.Discord.Guilds()
	guildlist.Guilds = make([]*pb.Guild, len(botguilds))

	for idx, guild := range botguilds {
		guildlist.Guilds[idx] = msgbuilder.Guild(guild)
	}

	return guildlist, nil
}

// GetGuild fetches a single Guild from state.
func (d *DiscordService) GetGuild(ctx context.Context, req *pb.IDQuery) (*pb.Guild, error) {
	g, err := d.Discord.Guild(req.GuildID)
	return msgbuilder.Guild(g), err
}

// GetGuildsByMember searches for guilds that include a certain member.
func (d *DiscordService) GetGuildsByMember(ctx context.Context, req *pb.IDQuery) (*pb.GuildList, error) {
	memberGuilds := &pb.GuildList{}
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
func (d *DiscordService) GetMember(ctx context.Context, req *pb.IDQuery) (*pb.Member, error) {
	member, err := d.Discord.GuildMember(req.MemberID, req.GuildID)
	return msgbuilder.Member(member), err
}

// UpdateMember
func (d *DiscordService) UpdateMember(ctx context.Context, req *pb.Member) (*pb.Member, error) {
	err := d.Discord.Session.GuildMemberEdit(req.GuildID, req.User.ID, req.Roles)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (d *DiscordService) GetGuildRoles(ctx context.Context, req *pb.IDQuery) (*pb.GuildRoles, error) {
	guild, err := d.Discord.Guild(req.GuildID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildRoles{
		ID:    guild.ID,
		Roles: msgbuilder.Roles(guild.Roles),
	}, nil
}

func (d *DiscordService) GetUser(ctx context.Context, req *pb.IDQuery) (*pb.User, error) {
	user, err := d.Discord.Session.User(req.MemberID)
	return msgbuilder.User(user), err
}

func (d *DiscordService) OwnUser(ctx context.Context, req *empty.Empty) (*pb.User, error) {
	return d.GetUser(ctx, &pb.IDQuery{MemberID: d.Discord.UserID()})
}
