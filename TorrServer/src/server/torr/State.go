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
	BytesWritten        int64
	BytesWrittenData    int64
	BytesRead           int64
	BytesReadData       int64
	BytesReadUsefulData int64
	ChunksWritten       int64
	ChunksRead          int64
	ChunksReadUseful    int64
	ChunksReadUnwanted  int64
	PiecesDirtiedGood   int64
	PiecesDirtiedBad    int64

	TotalPeers       int
	PendingPeers     int
	ActivePeers      int
	ConnectedSeeders int
	HalfOpenPeers    int

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
	Torrent     *torrent.Torrent `json:"-"`
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
	if info := ts.Torrent.Info(); info != nil {
		ts.TorrentSize = info.Length

		state := ts.Torrent.Stats()
		deltaDlBytes := state.BytesReadUsefulData.Int64() - ts.BytesReadUsefulData
		deltaUpBytes := state.BytesWrittenData.Int64() - ts.BytesWrittenData
		deltaTime := time.Since(ts.lastTimeSpeed).Seconds()

		ts.DownloadSpeed = float64(deltaDlBytes) / deltaTime
		ts.UploadSpeed = float64(deltaUpBytes) / deltaTime

		ts.BytesWritten = state.BytesWritten.Int64()
		ts.BytesWrittenData = state.BytesWrittenData.Int64()
		ts.BytesRead = state.BytesRead.Int64()
		ts.BytesReadData = state.BytesReadData.Int64()
		ts.BytesReadUsefulData = state.BytesReadUsefulData.Int64()
		ts.ChunksWritten = state.ChunksWritten.Int64()
		ts.ChunksRead = state.ChunksRead.Int64()
		ts.ChunksReadUseful = state.ChunksReadUseful.Int64()
		ts.ChunksReadUnwanted = state.ChunksReadUnwanted.Int64()
		ts.PiecesDirtiedGood = state.PiecesDirtiedGood.Int64()
		ts.PiecesDirtiedBad = state.PiecesDirtiedBad.Int64()

		ts.TotalPeers = state.TotalPeers
		ts.PendingPeers = state.PendingPeers
		ts.ActivePeers = state.ActivePeers
		ts.ConnectedSeeders = state.ConnectedSeeders
		ts.HalfOpenPeers = state.HalfOpenPeers
	}
	ts.Name = ts.Torrent.Name()
	ts.Hash = ts.Torrent.InfoHash().HexString()
	ts.LoadedSize = ts.Torrent.BytesCompleted()
	ts.lastTimeSpeed = time.Now()
}
