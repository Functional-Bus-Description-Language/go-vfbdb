package main

import (
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"

	"github.com/Functional-Bus-Description-Language/go-wbfbd/internal/args"
	"github.com/Functional-Bus-Description-Language/go-wbfbd/internal/vhdl"

	"fmt"
	"log"
)

func main() {
	log.SetFlags(0)

	cmdLineArgs := args.Parse()
	fmt.Println(cmdLineArgs)
	args.SetOutputPaths(cmdLineArgs)

	bus := fbdl.Compile(cmdLineArgs["global"]["main"])

	if _, ok := cmdLineArgs["vhdl"]; ok {
		vhdl.Generate(bus, cmdLineArgs)
	}
}
