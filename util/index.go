package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
