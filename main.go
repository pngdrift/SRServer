package main

import (
	"bufio"
	"os"
	"srserver/conf"
	"srserver/content"
	"srserver/logic"
	"srserver/logic/db"
	"strings"
)

func main() {
	contentServer := content.NewContentServer()
	go contentServer.Start()

	logicServer := logic.NewLogicServer()
	if !conf.TECHNICAL_WORKS {
		go logicServer.Start()
	}

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			switch input {
			case "shutdown", "stop", "restart":
				logicServer.Shutdown()
			case "update":
				contentServer.CreatePatchContainer()
				db.InitDb()
			}

		}
	}()

	select {}
}
