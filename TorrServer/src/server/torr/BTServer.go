package torr

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"server/settings"
	"server/torr/storage"
	"server/torr/storage/memcache"
	"server/utils"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/iplist"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/labstack/gommon/bytes"
	"golang.org/x/time/rate"
)

type BTServer struct {
	config *torrent.Config
	client *torrent.Client

	storage storage.Storage
	states  map[metainfo.Hash]*TorrentState

	mu sync.Mutex

	wmu      sync.Mutex
	watching bool
}

func NewBTS() *BTServer {
	bts := new(BTServer)
	bts.states = make(map[metainfo.Hash]*TorrentState)
	return bts
}

func (bt *BTServer) Connect() error {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	var err error
	bt.configure()
	bt.client, err = torrent.NewClient(bt.config)
	bt.states = make(map[metainfo.Hash]*TorrentState)
	return err
}

func (bt *BTServer) Disconnect() {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	if bt.client != nil {
		bt.client.Close()
		bt.client = nil
		utils.FreeOSMemGC()
	}
}

func (bt *BTServer) Reconnect() error {
	bt.Disconnect()
	return bt.Connect()
}

func (bt *BTServer) configure() {
	bt.storage = memcache.NewStorage(settings.Get().CacheSize)

	blocklist, _ := iplist.MMapPackedFile(filepath.Join(settings.Path, "blocklist"))

	userAgent := "uTorrent/3.4.9"
	peerID := "-UT3490-"

	bt.config = &torrent.Config{
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

		DefaultStorage: bt.storage,

		DhtStartingNodes: dht.GlobalBootstrapAddrs,
		ListenHost:       func(string) string { return "" },

		Bep20:         peerID,
		PeerID:        utils.PeerIDRandom(peerID),
		HTTPUserAgent: userAgent,

		EstablishedConnsPerTorrent: settings.Get().ConnectionsLimit,
		HalfOpenConnsPerTorrent:    int(float64(settings.Get().ConnectionsLimit) * 0.6),
		//TorrentPeersLowWater: 50,
		//TorrentPeersHighWater: 500,

		DisableIPv6: true,
	}

	if settings.Get().DownloadRateLimit > 0 {
		bt.config.DownloadRateLimiter = rate.NewLimiter(rate.Inf, settings.Get().DownloadRateLimit)
	}
	if settings.Get().UploadRateLimit > 0 {
		bt.config.UploadRateLimiter = rate.NewLimiter(rate.Inf, settings.Get().UploadRateLimit)
	}
}

func (bt *BTServer) Add(torrentLink string) (*settings.Torrent, error) {
	if bt.client == nil {
		return nil, errors.New("torrent client not started")
	}

	bt.mu.Lock()
	defer bt.mu.Unlock()

	mag, err := GetMagnet(torrentLink)
	if err != nil {
		return nil, err
	}

	return bt.add(mag)
}

func (bt *BTServer) Get(hashHex string) (*settings.Torrent, error) {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	return settings.LoadTorrentDB(hashHex)
}

func (bt *BTServer) Rem(hashHex string) error {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	return settings.RemoveTorrentDB(hashHex)
}

func (bt *BTServer) List() ([]*settings.Torrent, error) {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	return settings.LoadTorrentsDB()
}

func (bt *BTServer) BTState() *BTState {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	state := new(BTState)
	state.LocalPort = bt.client.LocalPort()
	state.PeerID = fmt.Sprintf("%x", bt.client.PeerID())
	state.BannedIPs = len(bt.client.BadPeerIPs())
	for _, dht := range bt.client.DhtServers() {
		state.DHTs = append(state.DHTs, dht)
	}
	for _, st := range bt.states {
		state.Torrents = append(state.Torrents, st)
	}
	return state
}

func (bt *BTServer) TorrentState(hashHex string) *TorrentState {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	hash := metainfo.NewHashFromHex(hashHex)
	if st, ok := bt.states[hash]; ok {
		return st
	}
	return nil
}

func (bt *BTServer) Clean(hashHex string) {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	if bt.storage != nil && hashHex == "" {
		bt.storage.Clean()
	} else if bt.client != nil && hashHex != "" {
		hash := metainfo.NewHashFromHex(hashHex)
		if tor, ok := bt.client.Torrent(hash); ok {
			delete(bt.states, hash)
			tor.Drop()
		}
	}
}

func (bt *BTServer) Preload(hashHex string, fileLink string) error {
	if settings.Get().PreloadBufferSize == 0 {
		return nil
	}

	tordb, err := bt.Get(hashHex)
	if err != nil {
		return err
	}

	var file *settings.File
	for _, f := range tordb.Files {
		if utils.FileToLink(f.Name) == fileLink {
			file = &f
			break
		}
	}

	if file == nil {
		return errors.New("File in torrent not found: " + hashHex + "/" + fileLink)
	}
	state, err := bt.getTorrent(tordb)
	if err != nil {
		return err
	}
	if !state.IsPreload {
		state.IsPreload = true
		pr := settings.Get().PreloadBufferSize
		ep := pr / state.PiecesLength
		if ep > 0 {
			state.PreloadLength = ep * state.PiecesLength
			state.torrent.DownloadPieces(0, int(ep))
			go bt.watcher()
			for {
				select {
				case <-state.torrent.Closed():
					return nil
				default:
				}
				state.expiredTime = time.Now().Add(time.Minute)
				state.PreloadSize = state.Filled
				fmt.Println("Preload:", bytes.Format(state.PreloadSize), "/", bytes.Format(state.PreloadLength), "Speed:", utils.Format(state.DownloadSpeed), "Peers:", state.ConnectedSeeders, "/", state.ActivePeers, ",", state.TotalPeers)
				if state.PreloadSize >= state.PreloadLength {
					return nil
				}
				time.Sleep(time.Second)
			}
		}
		state.IsPreload = false
	}
	return nil
}
