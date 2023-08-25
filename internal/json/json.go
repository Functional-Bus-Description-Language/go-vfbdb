package json

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/pkg"
)

func Generate(bus *fn.Block, pkgsConsts map[string]*pkg.Package, cmdLineArgs map[string]string) {
	err := os.MkdirAll(cmdLineArgs["-path"], os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	regFile, err := os.Create(path.Join(cmdLineArgs["-path"], "reg.json"))
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	byteArray, err := json.MarshalIndent(bus, "", "\t")
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	_, err = regFile.Write(byteArray)
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	err = regFile.Close()
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	constsFile, err := os.Create(path.Join(cmdLineArgs["-path"], "consts.json"))
	if err != nil {
		log.Fatalf("generate consts json: %v", err)
	}

	byteArray, err = json.MarshalIndent(pkgsConsts, "", "\t")
	if err != nil {
		log.Fatalf("generate consts json: %v", err)
	}

	_, err = constsFile.Write(byteArray)
	if err != nil {
		log.Fatalf("generate consts json: %v", err)
	}

	err = constsFile.Close()
	if err != nil {
		log.Fatalf("generate consts json: %v", err)
	}
}
