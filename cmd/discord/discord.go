package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	lru "github.com/hnlq715/golang-lru"
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

	memberCache, err := makeCache()
	if err != nil {
		klog.Fatal(err)
	}
	bot := setupBot()
	go startListener(bot, memberCache)

	go startGripkit(bot, memberCache)
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

func startGripkit(bot *discordclient.DiscordClient, cache *lru.ARCCache) {
	time.Sleep(2 * time.Second)
	grpcDiscord := rpcserver.NewDiscordService(bot, cache)

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
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(
					func(recovery interface{}) error {
						klog.Error("panic recovered: ", recovery)
						return nil
					},
				)),
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

func startListener(bot *discordclient.DiscordClient, cache *lru.ARCCache) {
	listener := &run.Listener{
		Bot:          bot,
		RootUsers:    strings.Split(os.Getenv("ROOT_USERS"), ","),
		BotWhitelist: strings.Split(os.Getenv("BOT_WHITELIST"), ","),
		YeetCache:    func() { cache.Purge() },
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

func makeCache() (*lru.ARCCache, error) {
	cacheTuningVar := os.Getenv("TUNING_MEMBER_CACHE_SIZE")
	if cacheTuningVar == "" {
		cacheTuningVar = "10000"
	}

	cacheTuning, err := strconv.Atoi(cacheTuningVar)
	if err != nil {
		klog.Warning("TUNING_MEMBER_CACHE_SIZE invalid, defauling to 10000")
		cacheTuning = 10000
	}

	memberCache, err := lru.NewARCWithExpire(cacheTuning, 2*time.Minute)
	if err != nil {
		klog.Fatal("Could not make memberCache")
	}

	return memberCache, err
}
