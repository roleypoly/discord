package permissions_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/roleypoly/discord/internal/permissions"
	"github.com/roleypoly/rpc/shared"
)

func TestCheckDiscordPermissions(t *testing.T) {
	testCases := []struct {
		desc       string
		role       *discordgo.Role
		permission int
		expect     bool
	}{
		{
			desc:       "No permissions doesn't have Administrator",
			role:       &discordgo.Role{Permissions: 0},
			permission: discordgo.PermissionAdministrator,
			expect:     false,
		},
		{
			desc:       "Admin permissions has Administrator",
			role:       &discordgo.Role{Permissions: discordgo.PermissionAdministrator},
			permission: discordgo.PermissionAdministrator,
			expect:     true,
		},
		{
			desc:       "Default permissions doesn't have Admin",
			role:       &discordgo.Role{Permissions: 104193601},
			permission: discordgo.PermissionAdministrator,
			expect:     false,
		},
		{
			desc:       "Default permissions doesn't have ManageRoles",
			role:       &discordgo.Role{Permissions: 104193601},
			permission: discordgo.PermissionManageRoles,
			expect:     false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if permissions.DiscordRoleHasPermission(tC.role, tC.permission) != tC.expect {
				t.Error("Incorrect")
			}
		})
	}
}

func TestCheckProtoPermissions(t *testing.T) {
	testCases := []struct {
		desc       string
		role       *shared.Role
		permission int
		expect     bool
	}{
		{
			desc:       "No permissions doesn't have Administrator",
			role:       &shared.Role{Permissions: 0},
			permission: discordgo.PermissionAdministrator,
			expect:     false,
		},
		{
			desc:       "Admin permissions has Administrator",
			role:       &shared.Role{Permissions: discordgo.PermissionAdministrator},
			permission: discordgo.PermissionAdministrator,
			expect:     true,
		},
		{
			desc:       "Default permissions doesn't have Admin",
			role:       &shared.Role{Permissions: 104193601},
			permission: discordgo.PermissionAdministrator,
			expect:     false,
		},
		{
			desc:       "Default permissions doesn't have ManageRoles",
			role:       &shared.Role{Permissions: 104193601},
			permission: discordgo.PermissionManageRoles,
			expect:     false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if permissions.ProtoRoleHasPermission(tC.role, tC.permission) != tC.expect {
				t.Error("Incorrect")
			}
		})
	}
}
