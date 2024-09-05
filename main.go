package main

import (
	"os"

	"github.com/geetchoubey/awsls/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(-1)
	}
}
