package torr

import (
	"time"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
)

type BTState struct {
	LocalPort int
	PeerID    string
	BannedIPs int
	DHTs      []*dht.Server

	Torrents []*TorrentState
}

type TorrentState struct {
	torrent.TorrentStats

	Name string
	Hash string

	IsPreload     bool
	LoadedSize    int64
	PreloadSize   int64
	PreloadLength int64

	lastTimeSpeed time.Time
	DownloadSpeed float64
	UploadSpeed   float64

	readers     int
	expiredTime time.Time
	torrent     *torrent.Torrent
}

func (ts *TorrentState) expired() bool {
	return ts.readers == 0 && ts.expiredTime.Before(time.Now())
}

func (ts *TorrentState) updateTorrentState() {
	if ts.torrent == nil {
		return
	}
	state := ts.torrent.Stats()

	deltaDlBytes := state.BytesReadUsefulData - ts.TorrentStats.BytesReadUsefulData
	deltaUpBytes := state.BytesWrittenData - ts.TorrentStats.BytesWrittenData
	deltaTime := time.Since(ts.lastTimeSpeed).Seconds()

	ts.DownloadSpeed = float64(deltaDlBytes) / deltaTime
	ts.UploadSpeed = float64(deltaUpBytes) / deltaTime
	ts.TorrentStats = state
	ts.Name = ts.torrent.Name()
	ts.Hash = ts.torrent.InfoHash().HexString()
	ts.LoadedSize = ts.torrent.BytesCompleted()
	ts.lastTimeSpeed = time.Now()
}
