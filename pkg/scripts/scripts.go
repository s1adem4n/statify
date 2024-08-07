package scripts

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"strings"
)

//go:embed tracker.min.js
var Tracker string

var TrackerHash string

func init() {
	hash := sha256.Sum256([]byte(Tracker))
	TrackerHash = hex.EncodeToString(hash[:])[:8]
}

func RenderTracker(serverAddress string) string {
	return strings.Replace(Tracker, "{{ .serverAddress }}", serverAddress, -1)
}
