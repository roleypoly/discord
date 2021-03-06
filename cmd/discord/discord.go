package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lampjaw/discordclient"
	"github.com/roleypoly/discord/cmd/discord/run"
	"github.com/roleypoly/discord/internal/version"
	"github.com/roleypoly/discord/rpcserver"
	"github.com/roleypoly/gripkit"
	proto "github.com/roleypoly/rpc/discord"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
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

var sharedSecret = os.Getenv("SHARED_SECRET")
var servicePort = os.Getenv("DISCORD_SVC_PORT")

func main() {
	klog.InitFlags(nil)
	klog.V(1).Info("Verbose on")
	klog.Infof(
		"Starting discord service.\n Build %s (%s) at %s",
		version.GitCommit,
		version.GitBranch,
		version.BuildDate,
	)

	defer awaitExit()

	bot := setupBot()

	go startGripkit(bot)
	go startListener(bot)
}

func setupBot() *discordclient.DiscordClient {
	client := discordclient.NewDiscordClient(discordConfig.BotToken, "", discordConfig.ClientID)
	client.AllowBots = true
	return client
}

func sharedSecretAuth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "shared")
	if err != nil {
		return nil, err
	}

	if token != sharedSecret {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}

	return ctx, nil
}

func startGripkit(bot *discordclient.DiscordClient) {
	grpcDiscord := rpcserver.NewDiscordService(bot)

	host, port, _ := net.SplitHostPort(os.Getenv("DISCORD_SVC_PORT"))
	healthzPort := host + ":1" + port

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
		gripkit.WithOptions(grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_auth.UnaryServerInterceptor(sharedSecretAuth),
			),
		)),
		gripkit.WithHealthz(&gripkit.HealthzOptions{
			UseDefault: true,
			Addr:       healthzPort,
		}),
	)

	proto.RegisterDiscordServer(gk.Server, grpcDiscord)

	klog.Info("gRPC server running on ", os.Getenv("DISCORD_SVC_PORT"))
	err := gk.Serve()
	if err != nil {
		klog.Exit("gRPC server failed to start.", err)
	}
}

func startListener(bot *discordclient.DiscordClient) {
	listener := &run.Listener{
		Bot:          bot,
		RootUsers:    strings.Split(os.Getenv("ROOT_USERS"), ","),
		BotWhitelist: strings.Split(os.Getenv("BOT_WHITELIST"), ","),
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
