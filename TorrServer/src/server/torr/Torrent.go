package torr

import (
	"fmt"
	"time"

	"server/settings"
	"server/utils"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func (bt *BTServer) add(magnet *metainfo.Magnet) (*settings.Torrent, error) {
	if len(magnet.Trackers) == 0 {
		magnet.Trackers = append(magnet.Trackers, utils.GetDefTrackers()...)
	}

	tor, _, err := bt.client.AddTorrentSpec(&torrent.TorrentSpec{
		Trackers:    [][]string{magnet.Trackers},
		DisplayName: magnet.DisplayName,
		InfoHash:    magnet.InfoHash,
	})

	if err != nil {
		return nil, err
	}

	fmt.Println("Adding", tor.Name())
	err = utils.GotInfo(tor, 20)
	if err != nil {
		return nil, err
	}
	go tor.Drop()
	torDb := new(settings.Torrent)
	torDb.Name = tor.Name()
	torDb.Hash = tor.InfoHash().HexString()
	torDb.Size = tor.Length()
	torDb.Magnet = magnet.String()
	torDb.Timestamp = time.Now().Unix()
	files := tor.Files()
	for _, f := range files {
		ff := settings.File{
			f.Path(),
			f.Length(),
			false,
		}
		torDb.Files = append(torDb.Files, ff)
	}
	err = settings.SaveTorrentDB(torDb)
	return torDb, err
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

			for hash, st := range bt.states {
				if st.expired() {
					bt.removeState(hash)
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
	bt.states[state.torrent.InfoHash()] = state
}

func (bt *BTServer) removeState(hash metainfo.Hash) {
	if st, ok := bt.states[hash]; ok {
		if st.torrent != nil {
			st.torrent.Drop()
		}
		delete(bt.states, hash)
	}
}

func (bt *BTServer) getTorrent(torrDB *settings.Torrent) (*TorrentState, error) {
	hash := metainfo.NewHashFromHex(torrDB.Hash)

	if st, ok := bt.states[hash]; ok {
		return st, nil
	}

	magnet, err := metainfo.ParseMagnetURI(torrDB.Magnet)
	if err != nil {
		return nil, err
	}

	switch settings.Get().RetrackersMode {
	case 1:
		magnet.Trackers = append(magnet.Trackers, utils.GetDefTrackers()...)
	case 2:
		magnet.Trackers = nil
	}

	tor, _, err := bt.client.AddTorrentSpec(&torrent.TorrentSpec{
		Trackers:    [][]string{magnet.Trackers},
		DisplayName: magnet.DisplayName,
		InfoHash:    magnet.InfoHash,
	})
	if err != nil {
		return nil, err
	}

	tor.SetMaxEstablishedConns(settings.Get().ConnectionsLimit)

	err = utils.GotInfo(tor, 60)
	if err != nil {
		return nil, err
	}

	state := new(TorrentState)
	state.torrent = tor
	state.TorrentSize = tor.Length()
	state.TorrentStats = tor.Stats()
	state.expiredTime = time.Now().Add(time.Minute * 5)
	bt.states[hash] = state

	return state, nil
}
