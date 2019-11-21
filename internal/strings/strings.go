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
)

type MentionResponseData struct {
	AppURL string
	GuildID string
}