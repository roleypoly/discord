package permissions

import (
	"github.com/bwmarrin/discordgo"
)

func RoleHasPermission(role *discordgo.Role, permission int) bool {
	return (role.Permissions & permission) == permission
}
