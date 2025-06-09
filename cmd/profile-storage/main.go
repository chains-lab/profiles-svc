package main

import (
	"os"

	"github.com/chains-lab/profile-storage/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
