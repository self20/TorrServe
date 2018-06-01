package main

import (
	"fmt"
	"os"
	"server"
	"server/settings"
)

func main() {
	path, _ := os.Getwd()
	if len(os.Args) == 2 {
		path = os.Args[1]
	}

	server.Start(path)
	settings.SaveSettings()

	fmt.Println(server.WaitServer())
}
