package main

import (
	"os"

	"github.com/br-lemes/golem/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
