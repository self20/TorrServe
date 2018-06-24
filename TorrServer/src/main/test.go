package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"server/settings"
	"server/torr/storage/memcache"
	"server/utils"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/iplist"
	"golang.org/x/time/rate"
)

func test() {
	config := torrent.NewDefaultClientConfig()

	config.EstablishedConnsPerTorrent = 150
	config.HalfOpenConnsPerTorrent = 97
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
		buf := make([]byte, 4*1024*1024)
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

func test2() {
	Magnet := []string{
		"magnet:?xt=urn:btih:FF21C570FCB1737CB18C9FABFC82D05B9F340BC1",
		//"magnet:?xt=urn:btih:164696AD3971A3AC81D989B9496CFBD065173F15",
	}

	userAgent := "uTorrent/3.4.9"
	peerID := "-UT3490-"

	storage := memcache.NewStorage(settings.Get().CacheSize)
	blocklist, _ := iplist.MMapPackedFile(filepath.Join(settings.Path, "blocklist"))

	config := &torrent.ClientConfig{
		DisableTCP:              settings.Get().DisableTCP,
		DisableUTP:              settings.Get().DisableUTP,
		NoDefaultPortForwarding: settings.Get().DisableUPNP,
		NoDHT:    settings.Get().DisableDHT,
		NoUpload: settings.Get().DisableUpload,

		EncryptionPolicy: torrent.EncryptionPolicy{
			DisableEncryption: settings.Get().Encryption == 1,
			ForceEncryption:   settings.Get().Encryption == 2,
		},
		DownloadRateLimiter: rate.NewLimiter(rate.Inf, 2<<16),
		UploadRateLimiter:   rate.NewLimiter(rate.Inf, 2<<16),

		IPBlocklist: blocklist,

		DefaultStorage: storage,

		DhtStartingNodes: dht.GlobalBootstrapAddrs,
		ListenHost:       func(network string) string { return "" },

		Bep20:         peerID,
		PeerID:        utils.PeerIDRandom(peerID),
		HTTPUserAgent: userAgent,

		EstablishedConnsPerTorrent: settings.Get().ConnectionsLimit,
		HalfOpenConnsPerTorrent:    int(float64(settings.Get().ConnectionsLimit) * 0.65),

		HandshakesTimeout: time.Second * 10,

		DisableIPv6: true,

		//Debug: true,
	}

	cl, err := torrent.NewClient(config)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}
	//for init and loading peers, dht, etc...
	time.Sleep(time.Second * 2)

	alltm := time.Duration(0)
	failed := 0
	count := 20
	for i := 0; i < count; i++ {
		fmt.Println("****** Connect", i, Magnet[0])
		t, err := cl.AddMagnet(Magnet[0])
		if err != nil {
			log.Fatalf("error adding magnet to client: %s", err)
		}
		tt := time.Now()
		err = utils.GotInfo(t, 20)
		if err != nil {
			failed++
			fmt.Println("****** Got info failed")
		}
		ts := time.Since(tt)
		alltm += ts
		fmt.Println("****** Got info", t.InfoHash().String(), ts)
		t.Drop()
		time.Sleep(time.Second)
	}
	fmt.Println("All time", alltm)
	fmt.Println("Failed", failed)
	fmt.Println("Average", alltm/time.Duration(count-failed))
}
