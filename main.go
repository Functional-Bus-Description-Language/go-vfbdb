package main

import (
	"github.com/Functional-Bus-Description-Language/go-wbfbd/internal/args"

	"fmt"
	"log"
)

func main() {
	log.SetFlags(0)

	cmdLineArgs := args.Parse()
	fmt.Println(cmdLineArgs)
}
