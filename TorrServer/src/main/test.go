package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

func test() {
	config := torrent.NewDefaultClientConfig()

	config.EstablishedConnsPerTorrent = 100
	config.HalfOpenConnsPerTorrent = 65
	config.DisableIPv6 = true
	config.NoDHT = true

	client, err := torrent.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	//Ubuntu
	t, err := client.AddMagnet("magnet:?xt=urn:btih:e4be9e4db876e3e3179778b03e906297be5c8dbe&dn=ubuntu-18.04-desktop-amd64.iso&tr=http%3a%2f%2ftorrent.ubuntu.com%3a6969%2fannounce&tr=http%3a%2f%2fipv6.torrent.ubuntu.com%3a6969%2fannounce")
	if err != nil {
		log.Fatal(err)
	}
	<-t.GotInfo()
	file := t.Files()[0]

	reader := file.NewReader()
	var wa sync.WaitGroup
	wa.Add(1)

	go func() {
		buf := make([]byte, t.Info().PieceLength)
		for {
			_, err := reader.Read(buf)
			if err != nil {
				break
			}
		}
		wa.Done()
	}()

	go func() {
		cl := t.Closed()
		for {
			select {
			case <-cl:
				return
			default:
				client.WriteStatus(os.Stdout)
			}
			time.Sleep(time.Second)
		}
	}()
	wa.Wait()
}
