package skeleton

import (
	"../util"
	"errors"
	"fmt"
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
	Subs        []SubSkeleton   `json:"subs"`
}

// command to run when build skeleton
type SkeletonCommand struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type SubSkeleton struct {
	SkeletonPath string `json:"skeletonPath"`
	TargetPath   string `json:"targetPath"`
	ContextPath  string `json:"contextPath"`
}

const (
	DEFAULT_TPL_DIR_NAME    = "stpl"
	DEFAULT_TARGET_DIR_NAME = "skel_target"
)

func BuildSkeleton(configFilePath string) {
	skeletonConfig := readSkeletonConfig(configFilePath)
	skeletonPath := getSkeletonClassPath(configFilePath, skeletonConfig.SkeletonPath)
	targetRoot := getTargetRootPath(filepath.Dir(configFilePath), skeletonConfig.TargetPath)

	RunSkeleton(skeletonPath, targetRoot, skeletonConfig.Context)
}

// @skeletonPath current skeleton class file path
// @targetRoot target root directory
// @context context map
func RunSkeleton(skeletonPath string, targetRoot string, context map[string]interface{}) {
	skeletonClass := readSkeletonClass(skeletonPath)

	// run before command
	if skeletonClass.Commands.Before != "" {
		util.Info("[gs-skeleton] run command :" + skeletonClass.Commands.Before)
		if err := runCommand(skeletonClass.Commands.Before, ""); err != nil {
			util.ExitWithError(err)
		}
	}

	// template dir
	skeletonTplPath := getSkeletonTemplateDirPath(skeletonPath, skeletonClass.TemplateDir)
	// template
	util.Info("[gs-skeleton] build template from dir:" + skeletonTplPath)
	buildErr := BuildTemplate(skeletonTplPath, targetRoot, context, TemplateOptions{
		Ignore: skeletonClass.Ignore,
	})

	if buildErr != nil {
		util.ErrorInfo("[gs-skeleton] fail to build template from :" + skeletonPath)
		util.ExitWithError(buildErr)
	}

	// build sub skeletons
	buildSubErr := buildSubSkeletons(skeletonClass.Subs, skeletonPath, targetRoot, context)

	if buildSubErr != nil {
		util.ExitWithError(buildSubErr)
	}

	// run after command
	if skeletonClass.Commands.After != "" {
		util.Info("[gs-skeleton] run command :" + skeletonClass.Commands.After)
		if err := runCommand(skeletonClass.Commands.After, targetRoot); err != nil {
			util.ExitWithError(err)
		}
	}
}

// TODO concurrent
func buildSubSkeletons(subs []SubSkeleton, skeletonPath string, targetRoot string, context map[string]interface{}) error {
	for _, sub := range subs {
		subSkeletonPath := getSkeletonClassPath(skeletonPath, sub.SkeletonPath)
		subTargetRoot := getTargetRootPath(targetRoot, sub.TargetPath)
		var subContext map[string]interface{}
		if sub.ContextPath != "" {
			nextContext, ok := util.GetByPath(context, sub.ContextPath)
			if !ok {
				return errors.New("missing sub context in path " + sub.ContextPath)
			}
			subContext, ok = nextContext.(map[string]interface{})
			if !ok {
				return errors.New("missing sub context in path " + sub.ContextPath)
			}
		}

		// run sub skeleton path
		RunSkeleton(subSkeletonPath, subTargetRoot, subContext)
	}

	return nil
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
func getTargetRootPath(currentDir, targetPath string) string {
	tarP := strings.TrimSpace(targetPath)
	if tarP == "" {
		tarP = DEFAULT_TARGET_DIR_NAME
	}
	return filepath.Join(currentDir, tarP)
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
//    - json file path
//    - directory with skl.json
//    - from global base directory
// @return the absolute file path for skeleton
// TODO download if not exists
func getSkeletonClassPath(configFilePath, skeletonPath string) string {
	return filepath.Join(filepath.Dir(configFilePath), skeletonPath)
}

// TODO check config
func readSkeletonConfig(configFilePath string) SkeletonConfig {
	util.Info("[gs-skeleton] read skeleton config file:" + configFilePath)
	var skeletonConfig SkeletonConfig
	util.ReadJsonWithPanic(configFilePath, &skeletonConfig, "[gs-skeleton] fail to read config content from :"+configFilePath)
	return skeletonConfig
}

// TODO check class
func readSkeletonClass(skeletonPath string) SkeletonClass {
	util.Info("[gs-skeleton] read skeleton class file:" + skeletonPath)
	var skeletonClass SkeletonClass
	util.ReadJsonWithPanic(skeletonPath, &skeletonClass, "[gs-skeleton] fail to read skeleton class from :"+skeletonPath)
	return skeletonClass
}
