# roley3 discord service

Microservice for the Discord bot portion of Roleypoly 3.0.

## Usage

Typical use case would be to run this through docker.

`docker run -it --rm -e DISCORD_BOT_TOKEN=... roleypoly/discord`

Otherwise, you may clone then build, or `go get github.com/roleypoly/discord`.

This bot doesn't do much with commands, as it's mostly a utility for the backend to pipe events into discord with.

## todo
- [x] bot: respond to mention (default case)
- [x] dgo: command router
- [x] pro: dockerize
- [x] pro: grpc
- [x] bot: respond to command (admin case)
- [ ] bot: dm auth challenge
- [ ] dgo: sharding

