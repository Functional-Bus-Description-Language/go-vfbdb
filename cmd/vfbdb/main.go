package main

import (
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"

	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/args"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/csync"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/json"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/python"
	"github.com/Functional-Bus-Description-Language/go-vfbdb/internal/vhdlwb3"

	"fmt"
	"log"
	"os"
)

var printDebug bool = false

type Logger struct{}

func (l Logger) Write(p []byte) (int, error) {
	print := true

	if len(p) > 4 && string(p)[:5] == "debug" {
		print = printDebug
	}

	if print {
		fmt.Fprint(os.Stderr, string(p))
	}

	return len(p), nil
}

func main() {
	logger := Logger{}
	log.SetOutput(logger)
	log.SetFlags(0)

	cmdLineArgs := args.Parse()
	args.SetOutputPaths(cmdLineArgs)

	if _, ok := cmdLineArgs["global"]["--debug"]; ok {
		printDebug = true
	}

	mainName := "Main"
	if _, ok := cmdLineArgs["global"]["-main"]; ok {
		mainName = cmdLineArgs["global"]["-main"]
	}
	addTimestamp := false
	if _, ok := cmdLineArgs["global"]["-add-timestamp"]; ok {
		addTimestamp = true
	}
	bus, pkgsConsts, err := fbdl.Compile(cmdLineArgs["global"]["main"], mainName, addTimestamp)
	if err != nil {
		log.Fatalf("compile: %v", err)
	}

	if _, ok := cmdLineArgs["json"]; ok {
		json.Generate(bus, pkgsConsts, cmdLineArgs["json"])
	}

	if _, ok := cmdLineArgs["c-sync"]; ok {
		csync.Generate(bus, pkgsConsts, cmdLineArgs["c-sync"])
	}

	if _, ok := cmdLineArgs["python"]; ok {
		python.Generate(bus, pkgsConsts, cmdLineArgs["python"])
	}

	if _, ok := cmdLineArgs["vhdl-wb3"]; ok {
		vhdlwb3.Generate(bus, pkgsConsts, cmdLineArgs["vhdl-wb3"])
	}

	if _, ok := cmdLineArgs["global"]["-fusesoc"]; ok {
		generateFuseSocCoreFile(cmdLineArgs["global"]["-fusesoc-vlnv"])
	}
}

func generateFuseSocCoreFile(fusesocVLNV string) {
	f, err := os.Create("main.core")
	if err != nil {
		log.Fatalf("generate FuseSoc .core file: %v", err)
	}
	defer f.Close()

	s := "CAPI=2:\n\n"
	s += fmt.Sprintf("name: %s\n\n", fusesocVLNV)
	s += "filesets:\n  vhdl:\n    depend: [mkru:vhdl-types:types]\n    file_type: vhdlSource-2008\n    logical_name: vfbdb\n    files:\n"

	for _, f := range vhdlwb3.GeneratedFiles {
		s += fmt.Sprintf("      - %s\n", f)
	}

	s += "\ntargets:\n  default:\n    filesets:\n      - vhdl"

	_, err = fmt.Fprint(f, s)
	if err != nil {
		log.Fatalf("generate FuseSoc.core file: %v", err)
	}
}
