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
	"github.com/roleypoly/gripkit"
	proto "github.com/roleypoly/rpc/discord"
)

var (
	token       = os.Getenv("DISCORD_BOT_TOKEN")
	rootUsers   = parseRoot(os.Getenv("ROOT_USERS"))
	svcPort     = os.Getenv("DISCORD_SVC_PORT")
	tlsCertPath = os.Getenv("TLS_CERT_PATH")
	tlsKeyPath  = os.Getenv("TLS_KEY_PATH")
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

	go startGripkit(discord)

	fmt.Println("roleypoly-discord: started grpc")

	syscallExit := make(chan os.Signal, 1)
	signal.Notify(
		syscallExit,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
		os.Kill,
	)
	<-syscallExit

	discord.Close()
}

func startGripkit(discord *discordgo.Session) {
	grpcDiscord := &DiscordService{
		Discord: discord,
	}

	gk := gripkit.Create(
		gripkit.WithHTTPOptions(gripkit.HTTPOptions{
			Addr:        os.Getenv("DISCORD_SVC_PORT"),
			TLSCertPath: os.Getenv("TLS_CERT_PATH"),
			TLSKeyPath:  os.Getenv("TLS_KEY_PATH"),
		}),
		gripkit.WithGrpcWeb(),
	)

	proto.RegisterDiscordServer(gk.Server, grpcDiscord)

	err := gk.Serve()
	if err != nil {
		log.Fatalln("grpc server failed to start", err)
	}
}
