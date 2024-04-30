package main

import (
	"fmt"
	"os"

	"github.com/noahc3/hatch-cli/src/loaders"
	"github.com/noahc3/hatch-cli/src/operations"
	"github.com/noahc3/hatch-cli/src/types"
	"github.com/noahc3/hatch-cli/src/utils"
)

var VALID_OPERATIONS = map[string]func(*types.Egg) error{
	"print": operations.Print,
}

func logf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		logf("Usage: hatch-cli <path> [<operation>...]\nExample: hatch-cli @github:pelican-eggs/eggs:game_eggs/minecraft/java/paper/egg-paper.json print\n\n")
		os.Exit(1)
	}

	tag := args[0]
	ops := args[1:]

	for _, operation := range ops {
		if !utils.KeysContains(VALID_OPERATIONS, operation) {
			logf("Invalid operation: %s\n", operation)
			os.Exit(1)
		}
	}

	egg, err := loaders.LoadEgg(tag)

	if err != nil {
		logf("Failed to load egg: %s\n", err)
		os.Exit(1)
	}

	for _, operation := range ops {
		err = VALID_OPERATIONS[operation](&egg)
		if err != nil {
			logf("Failed to run operation %s: %s\n", operation, err)
			os.Exit(1)
		}
	}
}
