package dm

import (
	"github.com/roleypoly/discord/internal/types"
	"regexp"
)

var Commands = []types.Command{
	{
		Matcher: regexp.MustCompile(`((log|sign) ?in|auth)`),
		Callback: func(message types.Message) string {
			return "ok one sec"
		},
	},
}
