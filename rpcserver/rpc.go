package rpcserver

import (
	"context"
	"os"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	lru "github.com/hashicorp/golang-lru"
	discordgobot "github.com/lampjaw/discordclient"
	"github.com/roleypoly/discord/internal/utils"
	"github.com/roleypoly/discord/msgbuilder"
	pbDiscord "github.com/roleypoly/rpc/discord"
	pbShared "github.com/roleypoly/rpc/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

// DiscordService is a gRPC implementation of rpc/discord.proto#Discord.
type DiscordService struct {
	pbDiscord.DiscordServer
	Discord *discordgobot.DiscordClient

	memberCache *lru.ARCCache
}

func NewDiscordService(discordClient *discordgobot.DiscordClient) *DiscordService {
	cacheTuningVar := os.Getenv("TUNING_MEMBER_CACHE_SIZE")
	if cacheTuningVar == "" {
		cacheTuningVar = "10000"
	}

	cacheTuning, err := strconv.Atoi(cacheTuningVar)
	if err != nil {
		klog.Warning("TUNING_MEMBER_CACHE_SIZE invalid, defauling to 10000")
		cacheTuning = 10000
	}

	memberCache, err := lru.NewARC(cacheTuning)
	if err != nil {
		klog.Fatal("Could not make memberCache")
	}

	return &DiscordService{
		Discord:     discordClient,
		memberCache: memberCache,
	}
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
		query := &pbShared.IDQuery{
			MemberID: req.MemberID,
			GuildID:  guild.ID,
		}
		mem, err := d.fetchMember(query, false)
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
	member, err := d.fetchMember(req, false)
	return msgbuilder.Member(member), err
}

// UpdateMember
func (d *DiscordService) UpdateMember(ctx context.Context, req *pbDiscord.Member) (*pbDiscord.Member, error) {
	err := d.Discord.Session.GuildMemberEdit(req.GuildID, req.User.ID, req.Roles)
	if err != nil {
		klog.Error("Failed to update roles: ", err)
		return nil, status.Error(codes.Internal, "Failed to update roles.")
	}

	return req, nil
}

// UpdateMemberRoles transactionally-ish updates roles with an add/remove action. Only makes one request, though.
func (d *DiscordService) UpdateMemberRoles(ctx context.Context, tx *pbDiscord.RoleTransaction) (*pbDiscord.RoleTransactionResult, error) {
	klog.Info("UpdateMemberRoles: got update for ", tx.Member, " using ", tx.Delta)
	member, err := d.fetchMember(tx.Member, true)
	if err != nil {
		klog.Error("UpdateMemberRoles: failed on fetch -- ", err)
		return nil, err
	}

	newRoles := member.Roles

	for _, delta := range tx.Delta {
		switch delta.Action {

		case pbDiscord.TxDelta_ADD:
			newRoles = append(newRoles, delta.Role)

		case pbDiscord.TxDelta_REMOVE:
			newRoles = utils.RemoveValueFromSlice(newRoles, delta.Role)
		}
	}

	err = d.Discord.Session.GuildMemberEdit(tx.Member.GuildID, tx.Member.MemberID, newRoles)
	if err != nil {
		klog.Error("UpdateMemberRoles: failed on edit -- ", err)
		return nil, grpc.Errorf(codes.Internal, "Role update failed.")
	}

	member.Roles = newRoles

	status := pbDiscord.RoleTransactionResult_DONE
	if d.isUpdateRatelimited(tx.Member.GuildID) {
		status = pbDiscord.RoleTransactionResult_QUEUED
	}

	d.memberCache.Remove(d.memberKey(tx.Member.GuildID, tx.Member.MemberID))

	return &pbDiscord.RoleTransactionResult{
		Member: msgbuilder.Member(member),
		Status: status,
	}, nil
}

func (d *DiscordService) GetGuildRoles(ctx context.Context, req *pbShared.IDQuery) (*pbShared.GuildRoles, error) {
	guild, err := d.Discord.Guild(req.GuildID)
	if err != nil {
		return nil, err
	}

	ownMember, err := d.ownMember(req.GuildID)
	if err != nil {
		return nil, err
	}

	roles := msgbuilder.Roles(guild.Roles)

	for _, role := range roles {
		role.Safety = calculateSafety(ownMember, guild, role)
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
