package cmd

import (
	"log"
	"os"
)

func outputRedirectedToFile() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		log.Printf("Cannot detect stdout type: %v", err.Error())
		return false
	}
	return info.Mode().IsRegular()
}
