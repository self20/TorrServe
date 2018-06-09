package torr

import (
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func (bt *BTServer) Watching(torr *torrent.Torrent) *TorrentState {
	state := new(TorrentState)
	state.Torrent = torr
	state.TorrentSize = torr.Length()
	state.TorrentStats = torr.Stats()
	state.expiredTime = time.Now().Add(time.Minute * 5)
	state.updateTorrentState()

	bt.states[torr.InfoHash()] = state
	bt.watcher()
	return state
}

func (bt *BTServer) watcher() {
	bt.wmu.Lock()
	defer bt.wmu.Unlock()
	if bt.watching {
		return
	}
	bt.watching = true
	go func() {
		for bt.watching {
			bt.wmu.Lock()

			for _, st := range bt.states {
				if st.expired() {
					bt.removeState(st)
				} else {
					st.updateTorrentState()
				}
			}

			bt.wmu.Unlock()
			time.Sleep(time.Second)
		}
		bt.watching = false
	}()
}

func (bt *BTServer) addState(state *TorrentState) {
	bt.states[state.Torrent.InfoHash()] = state
}

func (bt *BTServer) removeState(state *TorrentState) {
	hash := metainfo.NewHashFromHex(state.Hash)
	if st, ok := bt.states[hash]; ok {
		if st.Torrent != nil {
			st.Torrent.Drop()
		}
		delete(bt.states, hash)
	}
}
