package main

import (
	"./skeleton"
	"os"
	"path/filepath"
)

const DEFAULT_SKELTON_CONFIG = "gs.config.json"

// gs skelton.config
func main() {

	var configFilePath string
	cwd, cwdErr := os.Getwd()

	if cwdErr != nil {
		panic(cwdErr)
	}

	if len(os.Args) > 1 {
		configFilePath = filepath.Join(cwd, os.Args[1])
	} else {
		configFilePath = filepath.Join(cwd, DEFAULT_SKELTON_CONFIG)
	}

	skeleton.BuildSkelton(configFilePath)
}
