package vhdlwb3

import (
	"sync"
)

var GeneratedFiles []string
var genFilesMutex sync.Mutex

func addGeneratedFile(file string) {
	genFilesMutex.Lock()

	GeneratedFiles = append(GeneratedFiles, file)

	genFilesMutex.Unlock()
}
