package torrent

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"torrentserver/db"
	"torrentserver/settings"
	"torrentserver/storage/memcache"
	"torrentserver/utils"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/iplist"
	"github.com/anacrolix/torrent/metainfo"
	"golang.org/x/time/rate"

	"sync"
)

var (
	config  *torrent.Config
	client  *torrent.Client
	handler *Handler

	storage *memcache.Storage

	mutex sync.Mutex
)

func configure() {
	storage = memcache.NewStorage(settings.Get().CacheSize)

	blocklist, _ := iplist.MMapPackedFile(filepath.Join(settings.Get().SettingPath, "blocklist"))

	userAgent := "uTorrent/3.4.9"
	peerID := "-UT3490-"

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

		IPBlocklist: blocklist,

		DhtStartingNodes: dht.GlobalBootstrapAddrs,
		DefaultStorage:   storage,
		ListenHost:       func(string) string { return "" },

		Bep20:         peerID,
		PeerID:        utils.PeerIDRandom(peerID),
		HTTPUserAgent: userAgent,

		EstablishedConnsPerTorrent: settings.Get().ConnectionsLimit,
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

	handler = NewHandler()
	return nil
}

func Disconnect() {
	mutex.Lock()
	defer mutex.Unlock()
	if client != nil {
		handler.Close()
		client.Close()
		client = nil
		time.Sleep(time.Second * 3)
		runtime.GC()
		debug.FreeOSMemory()
	}
}

func Add(link string) (*db.Torrent, error) {
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

	defer tor.Drop()

	fmt.Println("Adding", tor.Name())
	err = GotInfo(tor)
	if err != nil {
		return nil, err
	}

	torDb := new(db.Torrent)
	torDb.Name = tor.Name()
	torDb.Hash = tor.InfoHash().HexString()
	torDb.Size = tor.Length()
	torDb.Magnet = mag
	torDb.Timestamp = time.Now().Unix()
	for _, f := range tor.Files() {
		ff := db.File{
			f.Path(),
			f.Length(),
			false,
		}
		torDb.Files = append(torDb.Files, ff)
	}
	err = db.SaveTorrentDB(torDb)
	return torDb, err
}

func Get(hashHex string) (*db.Torrent, error) {
	mutex.Lock()
	defer mutex.Unlock()

	return db.LoadTorrentDB(hashHex)
}

func Rem(hashHex string) error {
	mutex.Lock()
	defer mutex.Unlock()

	return db.RemoveTorrentDB(hashHex)
}

func List() ([]*db.Torrent, error) {
	mutex.Lock()
	defer mutex.Unlock()

	return db.LoadTorrentsDB()
}

func State(hashHex string) (*torrent.TorrentStats, error) {
	mutex.Lock()
	defer mutex.Unlock()
	hash := metainfo.NewHashFromHex(hashHex)
	if tor, ok := client.Torrent(hash); ok {
		st := tor.Stats()
		return &st, nil
	}
	return nil, errors.New("torrent not connected")
}

func CacheState() []memcache.CacheState {
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
