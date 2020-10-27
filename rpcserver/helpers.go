package rpcserver

import (
	"github.com/bwmarrin/discordgo"
	"github.com/roleypoly/discord/internal/permissions"
	"github.com/roleypoly/discord/internal/utils"
	"github.com/roleypoly/discord/msgbuilder"
	pbShared "github.com/roleypoly/rpc/shared"
	"k8s.io/klog"
)

// Fetch member by looking at cache, state, then REST. Returns nil if not present. Invalidate will skip cache.
func (d *DiscordService) fetchMember(req *pbShared.IDQuery, invalidate bool) (*discordgo.Member, error) {
	// guild, err := d.Discord.Guild(req.GuildID)
	// if err != nil {
	// 	return nil, err
	// }

	key := d.memberKey(req.GuildID, req.MemberID)
	if !invalidate && d.memberCache.Contains(key) {
		memberIntf, ok := d.memberCache.Get(key)
		if ok {
			if memberIntf == nil {
				return nil, nil
			}

			member, ok := memberIntf.(*discordgo.Member)
			if ok {
				return member, nil
			}
		}
	}

	member, err := d.Discord.GuildMember(req.MemberID, req.GuildID)
	if err != nil {
		if err != discordgo.ErrStateNotFound && err.Error() != "not found" {
			klog.Error("fetchMember (state) failed: ", req, " -- ", err)
			return nil, err
		}
	}

	if member == nil {
		// if guild.MemberCount > 5000 {
		member, err = d.Discord.Session.GuildMember(req.GuildID, req.MemberID)
		if err != nil && err.Error() != `HTTP 404 Not Found, {"message": "Unknown Member", "code": 10007}` {
			klog.Error("fetchMember (rest) failed: ", req, " -- ", err.Error())
		}
		// }

		if member == nil {
			d.memberCache.Add(key, nil)
			return nil, nil
		}
	}

	// This isn't set in every case, so let's make sure we do.
	member.GuildID = req.GuildID

	d.memberCache.Add(key, member)
	return member, nil
}

func (d *DiscordService) isUpdateRatelimited(guildID string) bool {
	bucketKey := discordgo.EndpointGuildMember(guildID, "")
	bucket := d.Discord.Session.Ratelimiter.GetBucket(bucketKey)
	if bucket.Remaining > 15 {
		if bucket.Remaining > 100 {
			klog.Info("threshold ratelimited in bucket: ", guildID, " -- ", bucket.Remaining, " remaining in queue")
		}

		return true
	}

	return false
}

func getRoleFromRoles(roleID string, roles []*discordgo.Role) *discordgo.Role {
	for _, role := range roles {
		if role.ID == roleID {
			return role
		}
	}

	return nil
}

func calculateSafety(targetMember *discordgo.Member, guild *discordgo.Guild, role *pbShared.Role) pbShared.Role_RoleSafety {
	if permissions.ProtoRoleHasPermission(role, discordgo.PermissionManageRoles) || permissions.ProtoRoleHasPermission(role, discordgo.PermissionAdministrator) {
		return pbShared.Role_dangerousPermissions
	}

	highestOwnRolePosition := 0
	for _, roleID := range targetMember.Roles {
		checkRole := getRoleFromRoles(roleID, guild.Roles)

		if checkRole.Position > highestOwnRolePosition {
			highestOwnRolePosition = checkRole.Position
		}
	}

	if role.Position > int32(highestOwnRolePosition) {
		return pbShared.Role_higherThanBot
	}

	return pbShared.Role_safe
}

func sanitizeRoles(targetMember *discordgo.Member, guild *discordgo.Guild, roles []string) []string {
	resultRoles := roles
	for _, roleID := range roles {
		role := getRoleFromRoles(roleID, guild.Roles)
		if role == nil {
			resultRoles = utils.RemoveValueFromSlice(resultRoles, roleID)
			continue
		}

		safety := calculateSafety(targetMember, guild, msgbuilder.Role(role))
		if safety != pbShared.Role_safe {
			resultRoles = utils.RemoveValueFromSlice(resultRoles, roleID)
			continue
		}
	}

	return resultRoles
}

func (d *DiscordService) ownMember(guildID string) (*discordgo.Member, error) {
	return d.fetchMember(&pbShared.IDQuery{
		MemberID: d.Discord.ClientID,
		GuildID:  guildID,
	}, false)
}

// type cacheKey string

func (d *DiscordService) memberKey(guildID, memberID string) string {
	return guildID + "-" + memberID
}
