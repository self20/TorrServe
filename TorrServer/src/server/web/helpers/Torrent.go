package helpers

import (
	"fmt"
	"io"
	"time"

	"server/settings"
	"server/torr"
	"server/utils"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func Add(bts *torr.BTServer, magnet *metainfo.Magnet, save bool) (*settings.Torrent, error) {
	if len(magnet.Trackers) == 0 {
		magnet.Trackers = append(magnet.Trackers, utils.GetDefTrackers()...)
	}

	fmt.Println("Adding torrent", magnet.String())
	torrState, err := bts.AddTorrent(magnet, 20)
	if err != nil {
		return nil, err
	}

	torDb := new(settings.Torrent)
	torDb.Name = torrState.Name
	torDb.Hash = torrState.Hash
	torDb.Size = torrState.TorrentSize
	torDb.Magnet = magnet.String()
	torDb.Timestamp = time.Now().Unix()
	files := torrState.Torrent.Files()
	for _, f := range files {
		ff := settings.File{
			f.Path(),
			f.Length(),
			false,
		}
		torDb.Files = append(torDb.Files, ff)
	}

	if save {
		err = settings.SaveTorrentDB(torDb)
	}
	if err != nil {
		return nil, err
	}

	return torDb, nil
}

func AddFile(bts *torr.BTServer, reader io.Reader) (*settings.Torrent, error) {
	info, err := metainfo.Load(reader)
	if err != nil {
		return nil, err
	}

	torrDb, err := settings.LoadTorrentDB(info.HashInfoBytes().String())
	if err != nil {
		return torrDb, nil
	}

	minfo, err := info.UnmarshalInfo()
	if err != nil {
		return nil, err
	}

	magnet := info.Magnet(minfo.Name, info.HashInfoBytes())

	torrDb, err = Add(bts, &magnet, true)
	return torrDb, err
}

func FindFile(fileLink string, torr *torrent.Torrent) *torrent.File {
	for _, f := range torr.Files() {
		if utils.FileToLink(f.Path()) == fileLink {
			return f
		}
	}
	return nil
}

//func getTorrent(torrDB *settings.Torrent) (*TorrentState, error) {
//	hash := metainfo.NewHashFromHex(torrDB.Hash)
//
//	if st, ok := bt.states[hash]; ok {
//		return st, nil
//	}
//
//	magnet, err := metainfo.ParseMagnetURI(torrDB.Magnet)
//	if err != nil {
//		return nil, err
//	}
//
//	switch settings.Get().RetrackersMode {
//	case 1:
//		magnet.Trackers = append(magnet.Trackers, utils.GetDefTrackers()...)
//	case 2:
//		magnet.Trackers = nil
//	}
//
//	tor, _, err := bt.client.AddTorrentSpec(&torrent.TorrentSpec{
//		Trackers:    [][]string{magnet.Trackers},
//		DisplayName: magnet.DisplayName,
//		InfoHash:    magnet.InfoHash,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	err = utils.GotInfo(tor, 60)
//	if err != nil {
//		go tor.Drop()
//		return nil, err
//	}
//
//	state := new(TorrentState)
//	state.torrent = tor
//	state.TorrentSize = tor.Length()
//	state.TorrentStats = tor.Stats()
//	state.expiredTime = time.Now().Add(time.Minute * 5)
//	bt.states[hash] = state
//
//	return state, nil
//}
