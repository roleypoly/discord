module github.com/roleypoly/discord

go 1.12

require (
	github.com/bwmarrin/discordgo v0.19.0
	github.com/improbable-eng/grpc-web v0.11.0
	github.com/joho/godotenv v1.3.0
	github.com/micro/go-micro v1.8.3 // indirect
	github.com/roleypoly/gripkit v0.0.0
	github.com/roleypoly/rpc v0.0.0-20190809193116-ff425f7c11c0
	google.golang.org/grpc v1.22.1
)

replace github.com/roleypoly/gripkit => ../gripkit
