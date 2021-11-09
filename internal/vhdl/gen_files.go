package vhdl

import (
	"sync"
)

var GeneratedFiles []string
var generatedFilesMutex sync.Mutex

func addGeneratedFile(file string) {
	generatedFilesMutex.Lock()

	GeneratedFiles = append(GeneratedFiles, file)

	generatedFilesMutex.Unlock()
}
