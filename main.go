package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/renevo/zombieutils/internal/command"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()

	if err := command.Execute(os.Args[1:]); err != nil {
		// kinda hacky, but only way to get around the double error message from cobra
		if strings.Contains(err.Error(), "unknown command") {
			return
		}

		fmt.Fprintf(os.Stderr, "failed to execute application %T: %s\n", err, err.Error())
		os.Exit(1)
	}
}
