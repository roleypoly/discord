package main

import (
	"context"

	"github.com/bwmarrin/discordgo"
	pb "github.com/roleypoly/rpc/discord"
)

type DiscordService struct {
	Discord *discordgo.Session
}

func (d *DiscordService) ListServers(ctx context.Context, req *pb.Empty) (*pb.GuildList, error) {
	guildlist := &pb.GuildList{}
	for _, guild := range d.Discord.State.Guilds {
		guildlist.Guilds = append(guildlist.Guilds, &pb.Guild{
			ID:          guild.ID,
			Name:        guild.Name,
			Icon:        guild.Icon,
			OwnerID:     guild.OwnerID,
			MemberCount: int32(guild.MemberCount),
			Splash:      guild.Splash,
		})
	}

	return guildlist, nil
}
