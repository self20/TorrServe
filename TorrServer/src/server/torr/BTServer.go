package torr

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"time"

	"server/settings"
	"server/torr/storage"
	"server/torr/storage/memcache"
	"server/torr/storage/state"
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

	states   map[metainfo.Hash]*TorrentState
	queueAdd map[metainfo.Hash]*TorrentState

	mu sync.Mutex

	wmu sync.Mutex
	qmu sync.Mutex

	watching bool
}

func NewBTS() *BTServer {
	bts := new(BTServer)
	bts.states = make(map[metainfo.Hash]*TorrentState)
	bts.queueAdd = make(map[metainfo.Hash]*TorrentState)
	return bts
}

func (bt *BTServer) Connect() error {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	var err error
	bt.configure()
	bt.client, err = torrent.NewClient(bt.config)
	bt.states = make(map[metainfo.Hash]*TorrentState)
	bt.queueAdd = make(map[metainfo.Hash]*TorrentState)
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
		DisableIPv6: true,

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
		ListenHost:       func(network string) string { return "" },

		Bep20:         peerID,
		PeerID:        utils.PeerIDRandom(peerID),
		HTTPUserAgent: userAgent,

		EstablishedConnsPerTorrent: settings.Get().ConnectionsLimit,
		HalfOpenConnsPerTorrent:    int(float64(settings.Get().ConnectionsLimit) * 0.65),
		//TorrentPeersLowWater: 50,
		//TorrentPeersHighWater: 500,

		//Debug: true,
	}

	if settings.Get().DownloadRateLimit > 0 {
		bt.config.DownloadRateLimiter = rate.NewLimiter(rate.Limit(settings.Get().DownloadRateLimit*1024), 1024)
	}
	if settings.Get().UploadRateLimit > 0 {
		bt.config.UploadRateLimiter = rate.NewLimiter(rate.Limit(settings.Get().UploadRateLimit*1024), 1024)
	}
}

func (bt *BTServer) addTorrent(magnet *metainfo.Magnet) (*TorrentState, *torrent.Torrent, error) {
	switch settings.Get().RetrackersMode {
	case 1:
		magnet.Trackers = append(magnet.Trackers, utils.GetDefTrackers()...)
	case 2:
		magnet.Trackers = nil
	case 3:
		magnet.Trackers = append(magnet.Trackers, utils.GetDefTrackers()...)
	}
	tor, _, err := bt.client.AddTorrentSpec(&torrent.TorrentSpec{
		Trackers:    [][]string{magnet.Trackers},
		DisplayName: magnet.DisplayName,
		InfoHash:    magnet.InfoHash,
	})

	if err != nil {
		return nil, nil, err
	}

	if st, ok := bt.states[magnet.InfoHash]; ok {
		return st, tor, nil
	}
	return nil, tor, nil
}

func (bt *BTServer) AddTorrentQueue(magnet *metainfo.Magnet, onAdd func(*TorrentState)) error {
	st, tor, err := bt.addTorrent(magnet)
	if err != nil {
		return err
	}

	if st != nil {
		onAdd(st)
		return nil
	}

	bt.addQueue(tor, onAdd)
	return nil
}

func (bt *BTServer) AddTorrent(magnet *metainfo.Magnet) (*TorrentState, error) {
	st, tor, err := bt.addTorrent(magnet)
	if err != nil {
		return nil, err
	}

	if st != nil {
		return st, nil
	}

	fmt.Println("Geting torrent info:", magnet.String())
	st = NewState(tor)
	st.IsGettingInfo = true
	bt.Watching(st)

	select {
	case <-tor.GotInfo():
		fmt.Println("Torrent received info:", st.Name)
		st.IsGettingInfo = false
		return st, nil
	case <-tor.Closed():
		return nil, errors.New("Torrent closed: " + st.Name)
	}
}

func (bt *BTServer) List() []*TorrentState {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	list := make([]*TorrentState, 0)
	for _, st := range bt.states {
		list = append(list, st)
	}

	for _, st := range bt.queueAdd {
		list = append(list, st)
	}
	return list
}

func (bt *BTServer) GetTorrent(hash metainfo.Hash) *TorrentState {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	if st, ok := bt.states[hash]; ok {
		return st
	}

	if st, ok := bt.queueAdd[hash]; ok {
		return st
	}
	return nil
}

func (bt *BTServer) RemoveTorrent(hashHex string) {
	bt.wmu.Lock()
	defer bt.wmu.Unlock()
	bt.removeState(hashHex)
	bt.removeQueue(hashHex)
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

func (bt *BTServer) CacheState(hash metainfo.Hash) *state.CacheState {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	state := bt.storage.GetStats(hash)
	return state
}

func (bt *BTServer) WriteState(w io.Writer) {
	bt.client.WriteStatus(w)
}

func (bt *BTServer) Clean(hashHex string) {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	if hashHex == "" {
		bt.storage.Clean()
	} else {
		hash := metainfo.NewHashFromHex(hashHex)
		if tor, ok := bt.client.Torrent(hash); ok {
			delete(bt.states, hash)
			tor.Drop()
		}
	}
}

func (bt *BTServer) Preload(hash metainfo.Hash, file *torrent.File) error {
	if settings.Get().PreloadBufferSize == 0 {
		return nil
	}

	state, ok := bt.states[hash]
	if !ok {
		return errors.New("File in Torrent not found: " + hash.HexString() + " | " + file.Path())
	}
	cState := bt.storage.GetStats(hash)

	if !state.IsPreload {
		state.IsPreload = true
		pr := settings.Get().PreloadBufferSize
		ep := int(pr / cState.PiecesLength)
		if ep > 0 {
			state.PreloadLength = int64(ep) * cState.PiecesLength
			state.Torrent.DownloadPieces(0, ep)
			go bt.watcher()
			cl := state.Torrent.Closed()
			var lastSize int64 = 0
			errCount := 0
			for {
				select {
				case <-cl:
					return nil
				default:
				}
				state.expiredTime = time.Now().Add(time.Minute)
				state.PreloadSize = state.LoadedSize
				fmt.Println("Preload:", bytes.Format(state.PreloadSize), "/", bytes.Format(state.PreloadLength), "Speed:", utils.Format(state.DownloadSpeed), "Peers:[", state.ConnectedSeeders, "]", state.ActivePeers, "/", state.TotalPeers)
				if state.PreloadSize >= state.PreloadLength {
					return nil
				}

				if lastSize == state.PreloadSize {
					errCount++
				} else {
					lastSize = state.PreloadSize
					errCount = 0
				}
				if errCount > 60 {
					return errors.New("long time no progress download")
				}

				time.Sleep(time.Second)
			}
		}
		state.IsPreload = false
	}
	return nil
}
