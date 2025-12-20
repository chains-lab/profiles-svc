package main

import (
	"os"

	"github.com/umisto/profiles-svc/cmd/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
