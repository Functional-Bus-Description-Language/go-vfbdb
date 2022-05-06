package args

import (
	"fmt"
	"os"
)

var helpMsg string = `Versatile Functional Bus Description Language compiler back-end written in Go.
Version: %s

Supported targets:
  - c-sync    C target with synchronous (blocking) interface functions,
  - python    Python target,
  - vhdl-wb3  VHDL target for Wishbone compilant with revision B.3.
To check valid flags and parameters for a given target type: 'vfbdb {target} -help'.

Usage:
  vfbdb [global flag or parameter] [{{target}} [target flag or parameter] ...] ... path/to/fbd/file/with/main/bus

  At least one target must be specified. The last argument is always a path
  to the fbd file containing a definition of the main bus, unless it is
  '-help' or '-version.'

Flags:
  -help     Display help.
  -version  Display version.
  -debug    Print debug messages.
  -fusesoc  Generate FuseSoc '.core' file.
            This flag rather should not be set manually.
            It is recommended to use vfbdb as a generator inside FuseSoc.
            All necessary files can be found in the 'FuseSoc' directory in the vfbdb repository.
  -times    Print compile and generate times.

Parameters:
  -fusesoc-vlnv  FuseSoc VLNV tag.
  -path  Path for target directories with output files.
         The default is 'vfbdb' directory in the current working directory.
`

func printHelp() {
	fmt.Printf(helpMsg, Version)
	os.Exit(0)
}

func printTargetHelp(target string) {
	switch target {
	case "c-sync":
		fmt.Printf(csyncHelpMsg)
	case "python":
		fmt.Printf(pythonHelpMsg)
	case "vhdl":
		fmt.Printf(vhdlHelpMsg)
	default:
		panic("should never happen")
	}

	os.Exit(0)
}

var csyncHelpMsg string = `Vfbdb help for C-Sync target.
C-Sync target is a C language target with synchronous (blocking) interface
functions.

Flags:
  -help        Display help.
  -no-asserts  Do not include asserts. Not yet implemented.

Parameters:
  -path  Path for output files.
`

var pythonHelpMsg string = `Vfbdb help for Python target.

Flags:
  -help        Display help.
  -no-asserts  Do not include asserts. Not yet implemented.

Parameters:
  -path  Path for output files.
`

var vhdlHelpMsg string = `Vfbdb help for VHDL target.

Flags:
  -help   Display help.
  -no-psl Do not include PSL assertions. Not yet implemented.

Parameters:
  -path  Path for output files.
`
