module github.com/roleypoly/discord

go 1.13

require (
	github.com/bwmarrin/discordgo v0.20.1
	github.com/golang/protobuf v1.3.3
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/improbable-eng/grpc-web v0.11.0
	github.com/joho/godotenv v1.3.0
	github.com/lampjaw/discordclient v0.0.0-20191203024148-a457503e9888
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/roleypoly/gripkit v0.0.0-20190819014327-7141453fff6a
	github.com/roleypoly/rpc v1.5.0
	golang.org/x/crypto v0.0.0-20191029031824-8986dd9e96cf // indirect
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c // indirect
	golang.org/x/sys v0.0.0-20191029155521-f43be2a4598c // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.27.0
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

//replace github.com/roleypoly/gripkit => ../gripkit
//replace github.com/roleypoly/rpc => ../rpc
//replace github.com/lampjaw/discordgobot => ../../lampjaw/discordgobot
