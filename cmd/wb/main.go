package main

import (
	"os"

	"lobster/internal/cli"
)

func main() {
	os.Exit(cli.New("workbuddy").Run(os.Args[1:]))
}
