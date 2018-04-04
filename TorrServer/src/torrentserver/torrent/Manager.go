package torrent

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"torrentserver/settings"
	"torrentserver/storage/memcache"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/iplist"
	"github.com/anacrolix/torrent/metainfo"
	"golang.org/x/time/rate"

	"sync"
)

var (
	config *torrent.Config
	client *torrent.Client

	storage *memcache.Storage

	mutex sync.Mutex
)

func configure() {
	storage = memcache.NewStorage(settings.Get().CacheSize)

	config = &torrent.Config{
		//Debug: true,

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

		DHTConfig: dht.ServerConfig{
			StartingNodes: dht.GlobalBootstrapAddrs,
		},
		DefaultStorage: storage,
		ListenAddr:     "0.0.0.0:0",
	}

	if settings.Get().DownloadRateLimit > 0 {
		config.DownloadRateLimiter = rate.NewLimiter(rate.Inf, settings.Get().DownloadRateLimit*1024)
	}
	if settings.Get().UploadRateLimit > 0 {
		config.UploadRateLimiter = rate.NewLimiter(rate.Inf, settings.Get().UploadRateLimit*1024)
	}
}

func Connect() error {
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	configure()
	client, err = torrent.NewClient(config)
	if err != nil {
		return err
	}

	blocklist, _ := iplist.MMapPackedFile(filepath.Join(settings.Get().SettingPath, "blocklist"))
	client.SetIPBlockList(blocklist)

	return loadTorrents()
}

func Disconnect() {
	mutex.Lock()
	defer mutex.Unlock()
	if client != nil {
		client.Close()
		client = nil
		time.Sleep(time.Second * 3)
		runtime.GC()
		debug.FreeOSMemory()
	}
}

func Add(link string) (*torrent.Torrent, error) {
	if client == nil {
		return nil, errors.New("list empty")
	}

	mutex.Lock()
	defer mutex.Unlock()

	mag, err := GetMagnet(link)
	if err != nil {
		return nil, err
	}
	tinfo, err := torrent.TorrentSpecFromMagnetURI(mag)
	if err != nil {
		return nil, err
	}

	tor, _, err := client.AddTorrentSpec(tinfo)
	if err != nil {
		return nil, err
	}
	fmt.Println("Adding", tor.Name())
	err = GotInfo(tor)
	if err != nil {
		return nil, err
	}

	addTime(tor)
	saveTorrents()
	return tor, nil
}

func Get(hashHex string) *torrent.Torrent {
	mutex.Lock()
	defer mutex.Unlock()

	hash := metainfo.NewHashFromHex(hashHex)
	tor, _ := client.Torrent(hash)
	GotInfo(tor)
	return tor
}

func Rem(hashHex string) {
	mutex.Lock()
	defer mutex.Unlock()

	hash := metainfo.NewHashFromHex(hashHex)
	if tor, ok := client.Torrent(hash); ok {
		fmt.Println("Remove:", tor.Name())
		tor.Drop()
		remTime(tor)
	}
	saveTorrents()
}

func List() []*torrent.Torrent {
	mutex.Lock()
	defer mutex.Unlock()

	torrs := client.Torrents()
	var wa sync.WaitGroup
	wa.Add(len(torrs))
	for _, t := range torrs {
		go func() {
			GotInfo(t)
			wa.Done()
		}()
	}
	wa.Wait()

	return client.Torrents()
}

func GetStates() []memcache.CacheState {
	if storage != nil {
		return storage.GetStats()
	}
	return nil
}

func CleanCache() {
	if storage != nil {
		storage.CleanCache()
	}
}
