package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/noahc3/hatch-cli/src/loaders"
	"github.com/noahc3/hatch-cli/src/log"
	"github.com/noahc3/hatch-cli/src/operations"
	"github.com/noahc3/hatch-cli/src/types"
	"github.com/noahc3/hatch-cli/src/utils"
)

type OptionalArgs struct {
	quiet      bool
	unattended bool
	help       bool
	outpath    string
}

var VALID_OPERATIONS = map[string]func(*types.Egg) error{
	"print":     operations.Print,
	"checkvars": operations.CheckVars,
	"install":   operations.Install,
}

func stop(format string, args ...interface{}) {
	log.Error(fmt.Sprintf("\nSTOP: %s", format), args...)
	os.Exit(1)
}

func extractOptionalArgs(args []string) (OptionalArgs, []string) {
	var optionalArgs OptionalArgs
	var newArgs []string

	skip := false

	for i, arg := range args {
		if skip {
			skip = false
			continue
		}

		if arg == "--quiet" || arg == "-q" {
			optionalArgs.quiet = true
		} else if arg == "--unattended" || arg == "-y" {
			optionalArgs.unattended = true
		} else if arg == "--help" {
			optionalArgs.help = true
		} else if arg == "--outpath" {
			if i+1 >= len(args) {
				stop("Missing argument for --outpath\n")
			}

			optionalArgs.outpath = args[i+1]
			skip = true
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	return optionalArgs, newArgs
}

func main() {
	optionalArgs, args := extractOptionalArgs(os.Args[1:])
	log.Quiet = optionalArgs.quiet
	utils.Unattended = optionalArgs.unattended

	if optionalArgs.help {
		log.Info("hatch-cli: A CLI tool for inspecting and executing Pterodactyl eggs without Pterodactyl or the Wings daemon.\n\n")
		log.Info("Usage: hatch-cli <tag> [<operation>...]\n")
		log.Info("  <tag> is a pointer to the Egg to use. This can be a file path, web URL, or a tag specifying remote provider information (more information below).\n")
		log.Info("  <operation> is a list of one or more valid operations:\n")
		log.Info("    print: Print the egg in JSON format.\n")
		log.Info("    checkvars: Load environment variables and check their values against the validation rules defiend by the Egg.\n")
		log.Info("  Optional arguments:\n")
		log.Info("    --quiet, -q: Suppress output.\n")
		log.Info("    --unattended, -y: Run in unattended mode (no prompts).\n")
		log.Info("    --help: Display this help message.\n")
		log.Info("\nSupported remote providers (square brackets indicate an optional parameter):\n")
		log.Info("  GitHub - @github:owner/repo[/ref]:path/to/file.json\n")
		log.Info("    Fetch the egg from a GitHub repository. 'ref' is optional and defaults to 'master'.\n")
		log.Info("    Example: @github:pelican-eggs/eggs:game_eggs/minecraft/java/paper/egg-paper.json\n")

		log.Info("\nExample: hatch-cli @github:pelican-eggs/eggs:game_eggs/minecraft/java/paper/egg-paper.json print\n")

		os.Exit(0)
	}

	if len(args) < 2 {
		log.Info("Usage: hatch-cli <tag> [<operation>...]\n")
		log.Info("Example: hatch-cli @github:pelican-eggs/eggs:game_eggs/minecraft/java/paper/egg-paper.json print\n")
		log.Info("Try hatch-cli --help for more information\n")
		os.Exit(1)
	}

	tag := args[0]
	ops := args[1:]

	for _, operation := range ops {
		if !utils.KeysContains(VALID_OPERATIONS, operation) {
			stop("Invalid operation: %s\n", operation)
		}
	}

	egg, err := loaders.LoadEgg(tag)

	if err != nil {
		stop("Failed to load egg: %s\n", err)
	}

	log.Nl()

	egg.CrackProps.InstallDirPath, err = filepath.Abs(optionalArgs.outpath)
	if err != nil {
		stop("Failed to get absolute path for output directory: %s\n", err)
	}

	for _, operation := range ops {
		err = VALID_OPERATIONS[operation](&egg)
		if err != nil {
			stop("Failed to run operation %s: %s\n", operation, err)
		}
	}
}
