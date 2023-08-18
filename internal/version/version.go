package version

import (
	"fmt"
	"strings"
)

var (
	Version   = "dev"
	GoVersion = "n/a"
	BuildTime = "n/a"
)

func BuildVersion() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Version %s", Version)
	if BuildTime != "n/a" {
		fmt.Fprintf(&sb, " built on %s", BuildTime)
	}
	if GoVersion != "n/a" {
		fmt.Fprintf(&sb, " (go v%s)", GoVersion)
	}
	return sb.String()
}
