package torr

import (
	"time"

	"github.com/anacrolix/torrent/metainfo"
)

func (bt *BTServer) Watching(state *TorrentState) {
	state.updateTorrentState()
	state.expiredTime = time.Now().Add(time.Minute * 5)
	bt.states[state.Torrent.InfoHash()] = state
	bt.watcher()
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
					bt.removeState(st.Hash)
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

func (bt *BTServer) removeState(hashHex string) {
	hash := metainfo.NewHashFromHex(hashHex)
	if st, ok := bt.states[hash]; ok {
		if st.Torrent != nil {
			st.Torrent.Drop()
		}
		delete(bt.states, hash)
	}
}
