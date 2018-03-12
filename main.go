package main

import (
	"./skeleton"
	"./util"
	"os"
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
		configFilePath = util.JoinPath(cwd, os.Args[1])
	} else {
		configFilePath = util.JoinPath(cwd, DEFAULT_SKELTON_CONFIG)
	}

	skeleton.BuildSkeleton(configFilePath)
}
