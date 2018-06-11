package helpers

import (
	"fmt"
	"time"

	"server/settings"
	"server/torr"
	"server/utils"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func Add(bts *torr.BTServer, magnet *metainfo.Magnet, save bool, timeout int) (*settings.Torrent, error) {
	if len(magnet.Trackers) == 0 {
		magnet.Trackers = append(magnet.Trackers, utils.GetDefTrackers()...)
	}

	fmt.Println("Adding torrent", magnet.String())
	torrState, err := bts.AddTorrent(magnet, timeout)
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

func FindFile(fileLink string, torr *torrent.Torrent) *torrent.File {
	for _, f := range torr.Files() {
		if utils.FileToLink(f.Path()) == fileLink {
			return f
		}
	}
	return nil
}
