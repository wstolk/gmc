package main

import (
	"os"

	"wstolk/gmc/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
