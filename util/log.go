package util

import (
	"log"
)

func Info(msg string) {
	log.Println("\033[33m " + msg)
}

func ErrorInfo(msg string) {
	log.Println("\033[31m " + msg)
}
