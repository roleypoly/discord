package strings

import (
	"text/template"
)

var (
	// MentionResponse is a template that requires strings AppURL and GuildID.
	MentionResponse = template.Must(
		template.New("MentionResponse").Parse(
			`:beginner: Assign your roles here! {{ .AppURL }}/s/{{ .GuildID }}`,
		),
	)

	RootStats = template.Must(
		template.New("RootStats").Parse(`🐈
**People Stats**
:blank: 🙎‍♀️ Total Users: {{ .Users }}
:blank: 👨‍👩‍👦‍👦 Total Guilds: {{ .Guilds }}
:blank: 🦺 Total Roles: {{ .Roles }}

**Bot Stats**
:blank: 🔩 Total Shards: {{ .Shards }}
:blank: ⚙️ Revision: {{ .GitCommit }} ({{ .GitBranch }})
:blank: ⏰ Built at {{ .BuildDate }}
`,
		),
	)
)

type MentionResponseData struct {
	AppURL  string
	GuildID string
}

type RootStatsData struct {
	Users     int
	Guilds    int
	Roles     int
	Shards    int
	GitCommit string
	GitBranch string
	BuildDate string
}
