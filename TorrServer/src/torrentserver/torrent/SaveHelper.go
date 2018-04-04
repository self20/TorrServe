package torrent

import (
	"fmt"
	"strconv"
	"time"

	"torrentserver/settings"

	"github.com/anacrolix/torrent"
)

func addTime(tor *torrent.Torrent) error {
	hash := tor.InfoHash().HexString()
	tm := fmt.Sprint(time.Now().Unix())
	return settings.SaveTorrentTime(hash, tm)
}

func GetTime(tor *torrent.Torrent) (int64, error) {
	hash := tor.InfoHash().HexString()
	tm, err := settings.GetTorrentTime(hash)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(tm, 10, 64)
}

func remTime(tor *torrent.Torrent) error {
	hash := tor.InfoHash().HexString()
	return settings.RemTorrentTime(hash)
}

func saveTorrents() error {
	err := settings.SaveTorrentsDB(asTorrentsDB(client.Torrents()))
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func loadTorrents() error {
	torrsDB, err := settings.ReadTorrentsDB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	torrents := toTorrents(torrsDB)
	fmt.Println("Load", len(torrents), "torrents")
	for _, ts := range torrents {
		t, _, err := client.AddTorrentSpec(ts)
		fmt.Println(t.Name())
		if err != nil {
			fmt.Println("Error load torrents:", err, ts.InfoHash.HexString())
		}
	}
	return nil
}

func asTorrentsDB(torrents []*torrent.Torrent) []string {
	ret := make([]string, len(torrents))
	for i, t := range torrents {
		ret[i] = Magnet(t)
	}
	return ret
}

func toTorrents(torrents []string) []*torrent.TorrentSpec {
	ret := make([]*torrent.TorrentSpec, len(torrents))
	for i, t := range torrents {
		ts, err := torrent.TorrentSpecFromMagnetURI(t)
		if err != nil {
			fmt.Println("Error load torrent from db:", err)
		} else {
			ret[i] = ts
		}
	}
	return ret
}
