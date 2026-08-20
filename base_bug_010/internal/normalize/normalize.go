package normalize

import (
	"regexp"
	"strings"
)

var num = regexp.MustCompile(`\b\d+(?:\.\d+)?\b`)
var uuid = regexp.MustCompile(`(?i)\b[0-9a-f]{8}-[0-9a-f-]{27,}\b`)
var secret = regexp.MustCompile(`(?i)(token|password|secret|api[_-]?key)=\S+`)

func Message(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = secret.ReplaceAllString(s, "$1=<redacted>")
	s = uuid.ReplaceAllString(s, "<uuid>")
	s = num.ReplaceAllString(s, "<n>")
	return strings.Join(strings.Fields(s), " ")
}
