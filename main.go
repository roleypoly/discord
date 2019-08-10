package main // import "github.com/roleypoly/discord"

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	_ "github.com/joho/godotenv/autoload"
	proto "github.com/roleypoly/rpc/discord"
	"google.golang.org/grpc"
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

func isGrpcWebHandling(wrappedGrpc *grpcweb.WrappedGrpcServer, r *http.Request) bool {
	return (wrappedGrpc.IsGrpcWebRequest(r))
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

	grpcDiscord := &DiscordService{
		Discord: discord,
	}

	grpcServer := grpc.NewServer()
	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)
	proto.RegisterDiscordServer(grpcServer, grpcDiscord)

	httpHandler := http.HandlerFunc(wrappedGrpc.ServeHTTP)
	err = http.ListenAndServeTLS(
		svcPort,
		tlsCertPath,
		tlsKeyPath,
		httpHandler,
	)
	if err != nil {
		log.Fatalln("grpc web failed to start", err)
	}

	syscallExit := make(chan os.Signal, 1)
	signal.Notify(syscallExit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-syscallExit

	discord.Close()
}
