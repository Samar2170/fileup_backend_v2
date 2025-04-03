package main

import (
	"fileupbackendv2/config"
	"fileupbackendv2/internal/dirManager"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "setup":
			checkSetup()
		case "server":
			RunServer()
		case "dir-utils":
			dirManager.UpdateDirsData()
			dirManager.UpdateUserDirsData()
		case "feature":
			log.Println("Feature under development")
		case "help", "usage", "-h":
			fmt.Println("Usage: fileupbackendv2 [setup|server|new-user|regenerate-key|check-key|feature]")
		default:
			log.Fatal("Invalid argument")
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			RunServer()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			StartCronServer()
			wg.Done()
		}()
		wg.Wait()
	}
}

func checkSetup() {
	// check if dir exists
	uploadsDir := config.UploadsDir
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadsDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	logsDir := config.LogsDir
	if logsDir == "" {
		logsDir = config.BaseDir + "/logs"
	}

	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err = os.MkdirAll(logsDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
