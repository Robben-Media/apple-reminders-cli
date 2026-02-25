package main

import (
	"os"

	"github.com/Robben-Media/apple-reminders-cli/cmd/appre/commands"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	commands.SetBuildInfo(version, buildTime)
	if err := commands.Execute(); err != nil {
		commands.WriteErrorJSON(err.Error())
		os.Exit(1)
	}
}
