package main

import (
	"os"

	"github.com/aftaab/trelay/cmd/trelay/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
