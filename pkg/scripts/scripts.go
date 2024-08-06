package scripts

import (
	_ "embed"
	"strings"
)

//go:embed tracker.min.js
var Tracker string

func RenderTracker(serverAddress string) string {
	return strings.Replace(Tracker, "{{ .serverAddress }}", serverAddress, -1)
}
