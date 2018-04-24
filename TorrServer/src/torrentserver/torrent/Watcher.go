package torrent

import (
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type TorrentStat struct {
	torrent.TorrentStats

	LastState     time.Time
	DownloadSpeed float64
	UploadSpeed   float64
}

type Watcher struct {
	hash metainfo.Hash

	isWatching bool
	mu         sync.Mutex

	state TorrentStat
}

//Create after connect torrent
func NewWatcher(hash metainfo.Hash) *Watcher {
	w := new(Watcher)
	w.hash = hash
	w.start()
	return w
}

func (w *Watcher) start() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.isWatching {
		return
	}
	if tor, ok := client.Torrent(w.hash); ok {
		w.isWatching = true
		w.state.TorrentStats = tor.Stats()
		w.state.LastState = time.Now()
	} else {
		return
	}

	go func() {
		for w.isWatching {
			time.Sleep(time.Second)
			if client == nil {
				continue
			}
			if tor, ok := client.Torrent(w.hash); ok {
				w.getSpeed(tor.Stats())
			} else {
				w.isWatching = false
			}
		}
	}()
}

func (w *Watcher) getSpeed(newState torrent.TorrentStats) {
	deltaDlBytes := newState.BytesReadData - w.state.BytesReadData
	deltaUpBytes := newState.BytesWrittenData - w.state.BytesWrittenData
	deltaTime := time.Since(w.state.LastState).Seconds()

	w.state.DownloadSpeed = float64(deltaDlBytes) / deltaTime
	w.state.UploadSpeed = float64(deltaUpBytes) / deltaTime
	w.state.TorrentStats = newState
	w.state.LastState = time.Now()
}

func (w *Watcher) GetState() *TorrentStat {
	return &w.state
}
