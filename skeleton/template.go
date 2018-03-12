package skeleton

import (
	"../util"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	DIRECTORY           = 0
	FILE                = 1
	OTHER_FILE_TYPE     = 10
	SKELETON_TPL_SUFFIX = ".stpl"
)

type fileStruct struct {
	Name     string
	Type     int
	Relative string // root
}

type TemplateOptions struct {
	Ignore []string
}

// copy template project to target directory.
// when meet filw while copying, try to parse template file with variable context
// 1. iteration to visit files
// 2. concurrency control
func BuildTemplate(tplDir string, // tpl directory path
	targetRoot string, // target directory
	context map[string]interface{}, // context map
	option TemplateOptions, //
) error {
	stack := []fileStruct{{
		Name:     filepath.Base(tplDir),
		Type:     DIRECTORY,
		Relative: ".",
	}}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1] //pop

		switch mode := top.Type; {
		case mode == DIRECTORY:
			curFilePath := getAbsPath(top.Relative, tplDir)
			// add new files into stack
			files, err := ioutil.ReadDir(curFilePath)

			if err != nil {
				return err
			}
			for _, file := range files {
				nextRelative := util.JoinPath(top.Relative, file.Name())
				// TODO using ignore
				doIgnore, igErr := shouldIgnore(nextRelative, option.Ignore)
				if igErr != nil {
					return igErr
				}
				if !doIgnore {
					stack = append(stack, fileStruct{
						Name:     file.Name(),
						Type:     getType(file),
						Relative: nextRelative,
					})
				}
			}
			// create directory if not exists
			cerr := createDir(getAbsPath(top.Relative, targetRoot))
			if cerr != nil {
				return cerr
			}
		case mode == FILE:
			// do the file template job
			ferr := handleFile(top, tplDir, targetRoot, context)
			if ferr != nil {
				return ferr
			}
		}
	}

	return nil
}

// support variable in fileName
func parseFileName(fileName string, context map[string]interface{}) (string, error) {
	tmpl, err := template.New("fileName").Parse(fileName)
	if err != nil {
		return fileName, nil
	}

	var tpl bytes.Buffer
	if err = tmpl.Execute(&tpl, context); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func shouldIgnore(relative string, ignores []string) (bool, error) {
	for _, ignore := range ignores {
		matched, err := filepath.Match(ignore, relative)
		if err != nil {
			return true, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

// copy normal file, but compile specific template file
// eg: a.js.stpl -> a.js
func handleFile(file fileStruct, tplDir string, targetRoot string, context map[string]interface{}) error {
	srcPath := getAbsPath(file.Relative, tplDir)

	isTplFile := strings.HasSuffix(file.Relative, SKELETON_TPL_SUFFIX)

	srcState, srcSErr := os.Stat(srcPath)
	if srcSErr != nil {
		return srcSErr
	}

	// get target file path
	targetFileRelative, ferr := parseFileName(file.Relative, context)
	if ferr != nil {
		return ferr
	}

	if isTplFile {
		targetFileRelative = targetFileRelative[0 : len(targetFileRelative)-len(SKELETON_TPL_SUFFIX)]
	}

	// detect target file
	targetPath := getAbsPath(targetFileRelative, targetRoot)

	// check exists
	tarfi, hasTargetErr := os.Stat(targetPath)
	if !os.IsNotExist(hasTargetErr) && tarfi.Mode().IsRegular() {
		return errors.New("Exists file with name " + targetPath)
	}

	if isTplFile {
		// parse template file
		tmpl, terr := template.ParseFiles(srcPath)
		if terr != nil {
			return terr
		}
		wf, werr := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, srcState.Mode().Perm())
		if werr != nil {
			return werr
		}

		return tmpl.Execute(wf, context)
	} else {
		return os.Link(srcPath, targetPath)
	}
}

func getType(file os.FileInfo) int {
	switch mode := file.Mode(); {
	case mode.IsRegular():
		return FILE
	case mode.IsDir():
		return DIRECTORY
	default:
		return OTHER_FILE_TYPE
	}
}

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	return nil
}

func getAbsPath(relative string, root string) string {
	return util.JoinPath(root, relative)
}
