package rpcserver

import (
	"testing"

	"github.com/roleypoly/discord/msgbuilder"
	"github.com/roleypoly/rpc/shared"
)

func TestCalculateSafety(t *testing.T) {

	testCases := []struct {
		desc   string
		role   *shared.Role
		target shared.Role_RoleSafety
	}{
		{
			desc:   "admin is dangerous",
			role:   msgbuilder.Role(testGuild.Roles[0]),
			target: shared.Role_dangerousPermissions,
		},
		{
			desc:   "mod is dangerous",
			role:   msgbuilder.Role(testGuild.Roles[3]),
			target: shared.Role_dangerousPermissions,
		},
		{
			desc:   "unpriv is higher",
			role:   msgbuilder.Role(testGuild.Roles[1]),
			target: shared.Role_higherThanBot,
		},
		{
			desc:   "color is safe",
			role:   msgbuilder.Role(testGuild.Roles[len(testGuild.Roles)-1]),
			target: shared.Role_safe,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			val := calculateSafety(testMember, testGuild, tC.role)
			if val != tC.target {
				t.Errorf("expected %s but got %s", tC.target, val)
			}
		})
	}
}

func findStringInSlice(haystack []string, needle string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}

func TestSanitizeRoles(t *testing.T) {

	input := []string{"color-blue", "color-red", "mod", "nonexist"}
	output := []string{"color-blue", "color-red"}

	result := sanitizeRoles(testMember, testGuild, input)

	if len(output) != len(result) {
		t.Errorf("result wrong length: %v", result)
	}

	for _, outputID := range output {
		if !findStringInSlice(result, outputID) {
			t.Errorf("%s missing in %v", outputID, result)
		}
	}

}
