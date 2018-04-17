package main

import (
	"fmt"
	"os"

	"torrentserver"
	"torrentserver/settings"
)

func main() {
	path, _ := os.Getwd()
	if len(os.Args) == 2 {
		path = os.Args[1]
	}

	torrentserver.Start(path)
	settings.SaveFile()

	fmt.Println(torrentserver.WaitServer())
}
