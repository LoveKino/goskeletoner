package skeleton

import (
	"../util"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
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

const (
	DEFAULT_TPL_DIR_NAME    = "stpl"
	DEFAULT_TARGET_DIR_NAME = "skel_target"
)

func BuildSkelton(configFilePath string) {
	skeletonConfig := readSkeletonConfig(configFilePath)
	skeletonPath := getSkeletonClassPath(configFilePath, skeletonConfig.SkeletonPath)
	targetRoot := getTargetRootPath(configFilePath, skeletonConfig.TargetPath)

	RunSkeleton(skeletonPath, targetRoot, skeletonConfig.Context)
}

func RunSkeleton(skeletonPath string, targetRoot string, context map[string]interface{}) {
	skeletonClass := readSkeletonClass(skeletonPath)

	skeletonTplPath := getSkeletonTemplateDirPath(skeletonPath, skeletonClass.TemplateDir)

	// run before command
	if skeletonClass.Commands.Before != "" {
		log.Print("[gs-skeleton] run command :" + skeletonClass.Commands.Before)
		if err := runCommand(skeletonClass.Commands.Before, ""); err != nil {
			panic(err)
		}
	}

	// template
	log.Print("[gs-skeleton] build template from dir:" + skeletonTplPath)
	buildErr := BuildTemplate(skeletonTplPath, targetRoot, context, TemplateOptions{
		Ignore: skeletonClass.Ignore,
	})

	if buildErr != nil {
		log.Print("[gs-skeleton] fail to build template from :" + skeletonPath)
		panic(buildErr)
	}

	// run after command
	if skeletonClass.Commands.After != "" {
		log.Print("[gs-skeleton] run command :" + skeletonClass.Commands.After)
		if err := runCommand(skeletonClass.Commands.After, targetRoot); err != nil {
			panic(err)
		}
	}

}

func runCommand(cmdStr string, dir string) error {
	cmd := exec.Command("bash", "-c", cmdStr)
	if dir != "" {
		cmd.Dir = dir
	}

	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)

	return nil
}

// using default if not provide
func getTargetRootPath(configFilePath, targetPath string) string {
	tarP := strings.TrimSpace(targetPath)
	if tarP == "" {
		tarP = DEFAULT_TARGET_DIR_NAME
	}
	return filepath.Join(filepath.Dir(configFilePath), tarP)
}

// using default tpl name if not provide
func getSkeletonTemplateDirPath(skeletonClassPath, templateDir string) string {
	tplDir := strings.TrimSpace(templateDir)
	if tplDir == "" {
		tplDir = DEFAULT_TPL_DIR_NAME
	}
	return filepath.Join(filepath.Dir(skeletonClassPath), tplDir)
}

// TODO default global directory
// require level
// @return the absolute file path for skeleton
// TODO download if not exists
func getSkeletonClassPath(configFilePath, skeletonPath string) string {
	return filepath.Join(filepath.Dir(configFilePath), skeletonPath)
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
