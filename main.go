package main

import (
	"github.com/mesosphere/bun/cmd"
)

// Populated by GoReleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
	cmd.CheckNewRelease()
	cmd.Execute()
}
