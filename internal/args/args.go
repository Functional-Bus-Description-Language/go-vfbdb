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
		case "-help":
			printHelp()
		case "-version":
			printVersion()
		default:
			if isValidTarget(arg) {
				log.Fatalf("missing main file, check 'wbfbd -help'")
			} else {
				log.Fatalf("'%s' is not valid target, check 'wbfbd -help'", arg)
			}
		}
	}

	args := map[string]map[string]string{"global": map[string]string{}}
	currentTarget := ""
	currentParam := ""
	expectArg := false
	inGlobalArgs := true

	for i, arg := range os.Args[1:] {
		if i == len(os.Args)-2 {
			if !strings.HasPrefix(arg, "-") {
				break
			}
		}

		if inGlobalArgs {
			if arg == "-help" {
				printHelp()
			} else if arg == "-version" {
				printVersion()
			} else if expectArg {
				args["global"][currentParam] = arg
				expectArg = false
			} else if arg == "-fusesoc" || arg == "-times" {
				args["global"][arg] = ""
			} else if arg == "-debug" {
				args["global"]["-debug"] = ""
			} else if arg == "-fusesoc-vlnv" || arg == "-path" {
				currentParam = arg
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
			args[currentTarget][currentParam] = arg
			expectArg = false
		} else if isValidTarget(arg) {
			currentTarget = arg
			args[arg] = map[string]string{}
		} else if !isValidParam(arg, currentTarget) &&
			!isValidFlag(arg, currentTarget) &&
			expectArg == false {
			log.Fatalf(
				"'%s' is not valid flag or parameter for '%s' target, "+
					"run 'wbfbd %[2]s -help' to see valid flags and parameters",
				arg, currentTarget,
			)
		} else if arg == "-help" {
			printTargetHelp(currentTarget)
		} else if isValidFlag(arg, currentTarget) {
			args[currentTarget][arg] = ""
		} else if isValidParam(arg, currentTarget) {
			currentParam = arg
			expectArg = true
		}
	}

	if expectArg {
		log.Fatalf("missing argument for '%s' parameter, target '%s'", currentParam, currentTarget)
	}

	args["global"]["main"] = os.Args[len(os.Args)-1]

	// Default values handling.
	if _, exists := args["global"]["-path"]; !exists {
		args["global"]["-path"] = "wbfbd"
	}

	if len(args) == 1 {
		fmt.Println("No target specified, run 'wbfbd -help' to check valid targets.")
		os.Exit(1)
	}

	return args
}

func SetOutputPaths(args map[string]map[string]string) {
	for target, v := range args {
		if target == "global" {
			continue
		}

		if _, exists := v["-path"]; exists {
			continue
		}

		args[target]["-path"] = args["global"]["-path"]
	}
}
