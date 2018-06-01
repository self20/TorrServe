package torr

import (
	"fmt"
	"strings"
	"time"

	"server/settings"
	"server/utils"

	"github.com/anacrolix/torrent/metainfo"
)

func (bt *BTServer) add(magnet string) (*settings.Torrent, error) {
	mag := magnet
	if !strings.Contains(magnet, "&tr=") {
		mag = utils.AddRetracker(magnet)
	}

	tor, err := bt.client.AddMagnet(mag)
	if err != nil {
		return nil, err
	}

	fmt.Println("Adding", tor.Name())
	err = utils.GotInfo(tor)
	if err != nil {
		return nil, err
	}

	torDb := new(settings.Torrent)
	torDb.Name = tor.Name()
	torDb.Hash = tor.InfoHash().HexString()
	torDb.Size = tor.Length()
	torDb.Magnet = magnet
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

	mag := torrDB.Magnet

	switch settings.Get().RetrackersMode {
	case 1:
		mag = utils.AddRetracker(mag)
	case 2:
		mag = utils.RemoveRetracker(mag)
	}

	tor, err := bt.client.AddMagnet(mag)
	if err != nil {
		return nil, err
	}
	tor.SetMaxEstablishedConns(settings.Get().ConnectionsLimit)
	err = utils.GotInfo(tor)
	if err != nil {
		return nil, err
	}

	state := new(TorrentState)
	state.torrent = tor
	state.TorrentStats = tor.Stats()
	state.expiredTime = time.Now().Add(time.Minute * 5)
	bt.states[hash] = state

	return state, nil
}
