package permissions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/roleypoly/rpc/shared"
)

func DiscordRoleHasPermission(role *discordgo.Role, permission int) bool {
	return (role.Permissions & permission) == permission
}

func ProtoRoleHasPermission(role *shared.Role, permission int) bool {
	return (role.Permissions & int64(permission)) == int64(permission)
}
