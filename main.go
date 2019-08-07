package main // import "github.com/roleypoly/discord"

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	proto "github.com/roleypoly/rpc/discord"
)

var (
	token     = os.Getenv("DISCORD_BOT_TOKEN")
	rootUsers = parseRoot(os.Getenv("ROOT_USERS"))
)

func parseRoot(s string) []string {
	return strings.Split(s, ",")
}

func main() {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}

	discord.AddHandler(messageHandler)

	err = discord.Open()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("roleypoly-discord: started bot")

	service := grpc.NewService(
		micro.Name("discord"),
	)

	service.Init()

	discordService := &DiscordService{
		Discord: discord,
	}
	proto.RegisterDiscordHandler(service.Server(), discordService)

	go func() {
		err := service.Run()
		if err != nil {
			log.Fatalf("roleypoly-discord: rpc service failed to start")
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
