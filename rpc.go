package main

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/roleypoly/discord/msgbuilder"
	pb "github.com/roleypoly/rpc/discord"
)

// DiscordService is a gRPC implementation of rpc/discord.proto#Discord.
type DiscordService struct {
	// pb.UnimplementedDiscordServer
	Discord *discordgo.Session
}

// ListGuilds lists every guild in state.
func (d *DiscordService) ListGuilds(ctx context.Context, req *empty.Empty) (*pb.GuildList, error) {
	guildlist := &pb.GuildList{}
	guildlist.Guilds = make([]*pb.Guild, len(d.Discord.State.Guilds))

	for idx, guild := range d.Discord.State.Guilds {
		guildlist.Guilds[idx] = msgbuilder.Guild(guild)
	}

	return guildlist, nil
}

// GetGuild fetches a single Guild from state.
func (d *DiscordService) GetGuild(ctx context.Context, req *pb.IDQuery) (*pb.Guild, error) {
	g, err := d.Discord.State.Guild(req.GuildID)
	return msgbuilder.Guild(g), err
}

// GetGuildsByMember searches for guilds that include a certain member.
func (d *DiscordService) GetGuildsByMember(ctx context.Context, req *pb.IDQuery) (*pb.GuildList, error) {
	memberGuilds := &pb.GuildList{}
	for _, guild := range d.Discord.State.Guilds {
		mem, err := d.Discord.State.Member(guild.ID, req.MemberID)
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
	member, err := d.Discord.State.Member(req.GuildID, req.MemberID)
	return msgbuilder.Member(member), err
}

// UpdateMember
func (d *DiscordService) UpdateMember(ctx context.Context, req *pb.Member) (*pb.Member, error) {
	err := d.Discord.GuildMemberEdit(req.GuildID, req.User.ID, req.Roles)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (d *DiscordService) GetGuildRoles(ctx context.Context, req *pb.IDQuery) (*pb.GuildRoles, error) {
	guild, err := d.Discord.State.Guild(req.GuildID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildRoles{
		ID:    guild.ID,
		Roles: msgbuilder.Roles(guild.Roles),
	}, nil
}

func (d *DiscordService) GetUser(ctx context.Context, req *pb.IDQuery) (*pb.User, error) {
	user, err := d.Discord.User(req.MemberID)
	return msgbuilder.User(user), err
}

func (d *DiscordService) OwnUser(ctx context.Context, req *pb.IDQuery) (*pb.User, error) {
	user := d.Discord.State.User
	return msgbuilder.User(user), nil
}
