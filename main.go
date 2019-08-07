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

	fmt.Println("Started roley3 discord.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
