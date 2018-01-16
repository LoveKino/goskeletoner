package skeleton

import (
	"../util"
	"log"
	"path/filepath"
)

// instance of skeleton
type SkeletonConfig struct {
	SkeletonPath string                 `json:"skeletonPath"`
	TargetPath   string                 `json:"targetPath"`
	Context      map[string]interface{} `json:"context"`
}

// skeleton class, used to define a skeleton
type SkeletonClass struct {
	TemplateDir string          `json:"templateDir"`
	Ignore      []string        `json:"ignore"`
	Commands    SkeletonCommand `json:"commands"`
}

// command to run when build skeleton
type SkeletonCommand struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

func BuildSkelton(configFilePath string) {
	skeletonConfig := readSkeletonConfig(configFilePath)

	skeletonPath := filepath.Join(filepath.Dir(configFilePath), skeletonConfig.SkeletonPath)
	skeletonClass := readSkeletonClass(skeletonPath)

	targetRoot := filepath.Join(filepath.Dir(configFilePath), skeletonConfig.TargetPath)
	skeletonTplPath := filepath.Join(filepath.Dir(skeletonPath), skeletonClass.TemplateDir)
	log.Print("[gs-skeleton] build template from dir:" + skeletonTplPath)
	buildErr := BuildTemplate(skeletonTplPath, targetRoot, skeletonConfig.Context, TemplateOptions{
		Ignore: skeletonClass.Ignore,
	})

	if buildErr != nil {
		log.Print("[gs-skeleton] fail to build template from :" + skeletonPath)
		panic(buildErr)
	}
}

// TODO check config
func readSkeletonConfig(configFilePath string) SkeletonConfig {
	log.Print("[gs-skeleton] read skeleton config file:" + configFilePath)
	var skeletonConfig SkeletonConfig
	util.ReadJsonWithPanic(configFilePath, &skeletonConfig, "[gs-skeleton] fail to read config content from :"+configFilePath)
	return skeletonConfig
}

// TODO check class
func readSkeletonClass(skeletonPath string) SkeletonClass {
	log.Print("[gs-skeleton] read skeleton class file:" + skeletonPath)
	var skeletonClass SkeletonClass
	util.ReadJsonWithPanic(skeletonPath, &skeletonClass, "[gs-skeleton] fail to read skeleton class from :"+skeletonPath)
	return skeletonClass
}
