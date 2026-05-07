package version

import (
	"fmt"
	"runtime"
)

// These variables are intended to be set via -ldflags at build time.
// Defaults keep local `go build`/`go run` usable.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func String() string {
	return fmt.Sprintf("devclean %s (%s, %s) %s/%s", Version, Commit, Date, runtime.GOOS, runtime.GOARCH)
}
