package main

import (
	"fmt"
	"os"

	"github.com/skanehira/ghf/cmd"
)

func main() {
	if err := initConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmd.Execute()
}
