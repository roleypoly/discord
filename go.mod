module github.com/roleypoly/discord

go 1.13

require (
	github.com/bwmarrin/discordgo v0.19.0
	github.com/golang/protobuf v1.3.2
	github.com/improbable-eng/grpc-web v0.11.0
	github.com/joho/godotenv v1.3.0
	github.com/nicksnyder/go-i18n/v2 v2.0.2
	github.com/roleypoly/gripkit v0.0.0-20190819014327-7141453fff6a
	github.com/roleypoly/rpc v0.0.0-20190921034711-ecdc7744d4c7
	golang.org/x/text v0.3.2
	google.golang.org/grpc v1.22.1
	gopkg.in/yaml.v2 v2.2.1
)

//replace github.com/roleypoly/gripkit => ../gripkit

//replace github.com/roleypoly/rpc => ../rpc
