module github.com/roleypoly/discord

go 1.12

require (
	github.com/bwmarrin/discordgo v0.19.0
	github.com/golang/protobuf v1.3.2
	github.com/improbable-eng/grpc-web v0.11.0
	github.com/joho/godotenv v1.3.0
	github.com/roleypoly/gripkit v0.0.0-20190819014327-7141453fff6a
	github.com/roleypoly/rpc v0.0.0-20190908183221-801f93147031
	google.golang.org/grpc v1.22.1
)

//replace github.com/roleypoly/gripkit => ../gripkit

//replace github.com/roleypoly/rpc => ../rpc
