package torr

import (
	"fmt"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func (bt *BTServer) addQueue(tor *torrent.Torrent, onAdd func(*TorrentState)) {
	go func() {
		fmt.Println("Geting torrent info:", tor.Name())
		st := NewState(tor)
		st.IsGettingInfo = true

		bt.qmu.Lock()
		bt.queueAdd[tor.InfoHash()] = st
		bt.qmu.Unlock()

		select {
		case <-tor.GotInfo():
			st.IsGettingInfo = false
			bt.Watching(st)
			fmt.Println("Torrent received info:", tor.Name())
			go onAdd(st)
		case <-tor.Closed():
			fmt.Println("Torrent closed:", tor.Name())
		}

		bt.qmu.Lock()
		delete(bt.queueAdd, tor.InfoHash())
		bt.qmu.Unlock()
	}()
}

func (bt *BTServer) removeQueue(hashHex string) {
	hash := metainfo.NewHashFromHex(hashHex)
	if st, ok := bt.queueAdd[hash]; ok {
		st.Torrent.Drop()
	}
}

func (bt *BTServer) listQueue() []*TorrentState {
	list := make([]*TorrentState, 0)
	for _, st := range bt.queueAdd {
		list = append(list, st)
	}
	return list
}
