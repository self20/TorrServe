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

	IsGettingInfo bool
	IsPreload     bool
	LoadedSize    int64
	PreloadSize   int64
	PreloadLength int64

	TorrentSize int64

	lastTimeSpeed time.Time
	DownloadSpeed float64
	UploadSpeed   float64

	readers     int
	expiredTime time.Time
	Torrent     *torrent.Torrent
}

func NewState(torr *torrent.Torrent) *TorrentState {
	state := new(TorrentState)
	state.Torrent = torr
	state.lastTimeSpeed = time.Now()
	state.updateTorrentState()
	return state
}

func (ts *TorrentState) expired() bool {
	return ts.readers == 0 && ts.expiredTime.Before(time.Now()) && !ts.IsGettingInfo
}

func (ts *TorrentState) Files() []*torrent.File {
	if ts.Torrent != nil && ts.Torrent.Info() != nil {
		return ts.Torrent.Files()
	}
	return nil
}

func (ts *TorrentState) updateTorrentState() {
	if ts.Torrent == nil {
		return
	}
	state := ts.Torrent.Stats()
	if info := ts.Torrent.Info(); info != nil {
		ts.TorrentSize = info.Length
	}
	deltaDlBytes := state.BytesReadUsefulData - ts.TorrentStats.BytesReadUsefulData
	deltaUpBytes := state.BytesWrittenData - ts.TorrentStats.BytesWrittenData
	deltaTime := time.Since(ts.lastTimeSpeed).Seconds()

	ts.DownloadSpeed = float64(deltaDlBytes) / deltaTime
	ts.UploadSpeed = float64(deltaUpBytes) / deltaTime
	ts.TorrentStats = state
	ts.Name = ts.Torrent.Name()
	ts.Hash = ts.Torrent.InfoHash().HexString()
	ts.LoadedSize = ts.Torrent.BytesCompleted()
	ts.lastTimeSpeed = time.Now()
}
