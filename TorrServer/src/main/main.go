package main

import (
	"fmt"
	"os"
	"server"
	"server/settings"
)

func main() {
	//test2()
	//return
	path, _ := os.Getwd()
	port := ""
	if len(os.Args) >= 2 {
		path = os.Args[1]
	}
	if len(os.Args) >= 3 {
		port = os.Args[2]
	}

	server.Start(path, port)
	settings.SaveSettings()

	fmt.Println(server.WaitServer())
}
