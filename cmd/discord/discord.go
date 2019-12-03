package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	_ "github.com/joho/godotenv/autoload"
	discordgobot "github.com/lampjaw/discordclient"
	"github.com/roleypoly/discord/cmd/discord/run"
	"github.com/roleypoly/discord/rpcserver"
	"github.com/roleypoly/gripkit"
	proto "github.com/roleypoly/rpc/discord"
)

type discordEnvConfig struct {
	ClientID     string
	ClientSecret string
	BotToken     string
}

var discordConfig = discordEnvConfig{
	ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
	ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
	BotToken:     os.Getenv("DISCORD_BOT_TOKEN"),
}

var sharedSecret = os.Getenv("DISCORD_SHARED_SECRET")
var servicePort = os.Getenv("DISCORD_SVC_PORT")

func main() {
	defer awaitExit()

	bot := setupBot()

	go startGripkit(bot)
	go startListener(bot)
}

func setupBot() *discordgobot.DiscordClient {
	client := discordgobot.NewDiscordClient(discordConfig.BotToken, "", discordConfig.ClientID)

	return client
}

func startGripkit(bot *discordgobot.DiscordClient) {
	grpcDiscord := &rpcserver.DiscordService{
		Discord: bot,
	}

	gk := gripkit.Create(
		gripkit.WithHTTPOptions(gripkit.HTTPOptions{
			Addr:        os.Getenv("DISCORD_SVC_PORT"),
			TLSCertPath: os.Getenv("TLS_CERT_PATH"),
			TLSKeyPath:  os.Getenv("TLS_KEY_PATH"),
		}),
		gripkit.WithGrpcWeb(
			grpcweb.WithOriginFunc(func(o string) bool { return true }),
		),
		gripkit.WithDebug(),
	)

	proto.RegisterDiscordServer(gk.Server, grpcDiscord)

	log.Println("gRPC server running")
	err := gk.Serve()
	if err != nil {
		log.Fatalln("grpc server failed to start", err)
	}
}

func startListener(bot *discordgobot.DiscordClient) {
	listener := &run.Listener{
		Bot: bot,
	}

	listener.Run()
}

func awaitExit() {
	syscallExit := make(chan os.Signal, 1)
	signal.Notify(
		syscallExit,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
		os.Kill,
	)
	<-syscallExit
}
