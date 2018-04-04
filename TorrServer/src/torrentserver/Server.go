package torrentserver

import (
	"fmt"
	"time"

	"torrentserver/server"
	"torrentserver/settings"
)

func Start(setpath string) {
	if setpath != "" {
		err := settings.LoadFile(setpath)
		if err != nil {
			fmt.Println("Error load settings on start server:", setpath)
		}
	}
	server.Start()
	time.Sleep(time.Second)
}

func WaitServer() string {
	err := server.Wait()
	if err != nil {
		return err.Error()
	}
	return ""
}

func Stop() {
	go server.Stop()
}
