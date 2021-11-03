// Custom package for command line arguments parsing.
package args

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const Version string = "0.0.0"

func printVersion() {
	fmt.Println(Version)
	os.Exit(0)
}

func Parse() map[string]map[string]string {
	if len(os.Args) == 1 {
		printHelp()
	}

	if len(os.Args) == 2 {
		arg := os.Args[1]
		switch arg {
		case "-h", "--help":
			printHelp()
		case "-v", "--version":
			printVersion()
		default:
			if isValidTarget(arg) {
				log.Fatalf("missing main file, check 'wbfbd --help'")
			} else {
				log.Fatalf("'%s' is not valid target, check 'wbfbd --help'", arg)
			}
		}
	}

	args := map[string]map[string]string{"global": map[string]string{}}
	currentTarget := ""
	currentOption := ""
	expectArg := false
	inGlobalArgs := true

	for i, arg := range os.Args[1:] {
		if i == len(os.Args)-2 {
			if !strings.HasPrefix(arg, "-") {
				break
			}
		}

		if inGlobalArgs {
			if arg == "-h" || arg == "--help" {
				printHelp()
			} else if arg == "-v" || arg == "--version" {
				printVersion()
			} else if expectArg {
				args["global"][currentOption] = arg
				expectArg = false
			} else if arg == "--fusesoc" || arg == "--times" {
				args["global"][arg] = ""
			} else if arg == "--fusesoc-vlnv" || arg == "--path" {
				currentOption = arg
				expectArg = true
			} else if !strings.HasPrefix(arg, "-") {
				inGlobalArgs = false
				if !isValidTarget(arg) {
					log.Fatalf("'%s' is not valid target", arg)
				}
				currentTarget = arg
				args[arg] = map[string]string{}
			}

			continue
		}

		if expectArg {
			args[currentTarget][currentOption] = arg
			expectArg = false
		} else if isValidTarget(arg) {
			currentTarget = arg
			args[arg] = map[string]string{}
		} else if !isValidOption(arg, currentTarget) &&
			!isValidFlag(arg, currentTarget) &&
			expectArg == false {
			log.Fatalf(
				"'%s' is not valid flag or option for '%s' target, "+
					"run 'wbfbd %[2]s --help' to see valid flags and options",
				arg, currentTarget,
			)
		} else if arg == "-h" || arg == "--help" {
			printTargetHelp(currentTarget)
		} else if isValidFlag(arg, currentTarget) {
			args[currentTarget][arg] = ""
		} else if isValidOption(arg, currentTarget) {
			currentOption = arg
			expectArg = true
		}
	}

	if expectArg {
		log.Fatalf("missing argument for '%s' option, target '%s'", currentOption, currentTarget)
	}

	args["global"]["main"] = os.Args[len(os.Args)-1]

	// Default values handling.
	if _, exists := args["global"]["--path"]; !exists {
		args["global"]["--path"] = "wbfbd"
	}

	return args
}
