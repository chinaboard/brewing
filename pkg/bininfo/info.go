package bininfo

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	GitStatus      = "unknown"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
)

func StringifySingleLine() string {
	return fmt.Sprintf("GitStatus=%s. BuildTime=%s. GoVersion=%s. Runtime=%s/%s.",
		GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func StringifyMultiLine() string {
	return fmt.Sprintf("\nGitStatus=%s\nBuildTime=%s\nGoVersion=%s\nRuntime=%s/%s\n",
		GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func beauty() {
	if GitStatus == "" {
		GitStatus = "cleanly"
	} else {
		GitStatus = strings.Replace(strings.Replace(GitStatus, "\r\n", " |", -1), "\n", " |", -1)
	}
}

func init() {
	beauty()
}
