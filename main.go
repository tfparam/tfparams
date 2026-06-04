package main

import (
	"fmt"
	"os"

	"github.com/tfkit/tfparams/cmd"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "tfparams:", err)
		os.Exit(1)
	}
}
