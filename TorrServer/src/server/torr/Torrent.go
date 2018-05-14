package torr

import (
	"fmt"
	"sort"
	"time"

	"server/settings"
	"server/utils"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func (bt *BTServer) add(magnet string) (*settings.Torrent, error) {
	tinfo, err := torrent.TorrentSpecFromMagnetURI(magnet)
	if err != nil {
		return nil, err
	}

	tor, _, err := bt.client.AddTorrentSpec(tinfo)
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
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path() < files[j].Path()
	})
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
					st.updateCacheState(bt.storage.GetStats(st.torrent.InfoHash()))
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

	tor, err := bt.client.AddMagnet(torrDB.Magnet)
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
	state.CacheState = bt.storage.GetStats(hash)
	state.expiredTime = time.Now().Add(time.Minute * 5)
	bt.states[hash] = state
	go bt.watcher()

	return state, nil
}
