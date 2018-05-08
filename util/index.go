package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO support remote url
func ReadJson(path string, inst interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, inst)
}

func ReadJsonWithPanic(path string, inst interface{}, errMsg string) {
	err := ReadJson(path, inst)
	if err != nil {
		log.Print(errMsg)
		panic(err)
	}
}

func ExitWithError(err error) {
	ErrorInfo(err.Error())
	os.Exit(1)
}

func GetByPath(context map[string]interface{}, jsonPath string) (interface{}, bool) {
	jsonPath = strings.TrimSpace(jsonPath)
	if jsonPath == "." {
		return context, true
	}

	var curContext interface{}
	curContext = context
	parts := strings.Split(jsonPath, ".")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			mapCtx, cok := curContext.(map[string]interface{})
			if !cok {
				return nil, false
			}
			next, ok := mapCtx[part]
			if !ok {
				return nil, false
			}
			curContext = next
		}
	}

	return curContext, true
}

func JoinPath(path1, path2 string) string {
	if path2[0] == '/' {
		return path2
	}

	return filepath.Join(path1, path2)
}
