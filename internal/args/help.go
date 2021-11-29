package args

import (
	"fmt"
	"os"
)

var helpMsg string = `Functional Bus Description Language compiler back-end written in Go.
Version: %s

Supported targets: python, vhdl
To check valid flags and options for a given target type: 'wbfbd {target} --help'

Usage:
  wbfbd [{{target}} [target flag or option] ...] ... path/to/fbd/file/with/main/bus

  At least one target must be specified. The last argument is always a path
  to the fbd file containing a definition of the main bus, unless it is
  '-h', '--help', '-v' or '--version.'

Flags:
  -h, --help     Display help.
  -v, --version  Display version.
  -d, --debug    Print debug messages.
  --fusesoc  Generate FuseSoc '.core' file.
             This flag rather should not be set manually.
             It is recommended to use wbfbd as a generator inside FuseSoc.
             All necessary files can be found in the 'FuseSoc' directory in the wbfbd repository.
  --times  Print compile and generate times.

Options:
  --fusesoc-vlnv  FuseSoc VLNV tag.
  --path  Path for target directories with output files.
          The default is 'wbfbd' directory in the current working directory.
`

func printHelp() {
	fmt.Printf(helpMsg, Version)
	os.Exit(0)
}

func printTargetHelp(target string) {
	switch target {
	case "python":
		printPythonHelp()
	case "vhdl":
		printVHDLHelp()
	default:
		panic("should never happen")
	}
}

var pythonHelpMsg string = `wbfbd help for Python target

Flags:
  -h, --help    Display help.
  --no-asserts  Do not include asserts. Not yet implemented.

Options:
  --path  Path for output files.
`

func printPythonHelp() {
	fmt.Printf(pythonHelpMsg)
	os.Exit(0)
}

var vhdlHelpMsg string = `wbfbd help for VHDL target

Flags:
  -h, --help  Display help.
  --no-psl    Do not include PSL assertions. Not yet implemented.

Options:
  --path  Path for output files.
`

func printVHDLHelp() {
	fmt.Printf(vhdlHelpMsg)
	os.Exit(0)
}
