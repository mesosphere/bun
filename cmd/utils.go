package cmd

import (
	"log"
	"os"
)

func outputRedirectedToFile() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		log.Println("Cannot detect stdout type")
	}
	return info.Mode().IsRegular()
}
