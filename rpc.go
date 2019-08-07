package main

import (
	"context"

	"github.com/bwmarrin/discordgo"
	pb "github.com/roleypoly/rpc/discord"
)

type DiscordService struct {
	Discord *discordgo.Session
}

func (d *DiscordService) RootGetAllServers(ctx context.Context, req *pb.Empty, rsp *pb.ServerSlugPayload) error {
	rsp.Servers = []*pb.ServerSlug{
		{
			Id:   "111",
			Name: "test",
		},
	}

	return nil
}
