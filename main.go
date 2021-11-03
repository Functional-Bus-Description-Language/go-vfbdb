package main

import (
	"github.com/Functional-Bus-Description-Language/go-wbfbd/internal/args"

	"fmt"
)

func main() {
	cmdLineArgs := args.Parse()
	fmt.Println(cmdLineArgs)
}
