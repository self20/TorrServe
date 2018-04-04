package torrent

import (
	"errors"
	"sort"
	"time"

	"github.com/anacrolix/torrent"
)

func GotInfo(t *torrent.Torrent) error {
	select {
	case <-t.GotInfo():
		return nil
	case <-time.Tick(time.Second * 15):
		return errors.New("timeout load torrent info")
	}
}

func Magnet(t *torrent.Torrent) string {
	inf := t.Metainfo()
	return (&inf).Magnet(t.Name(), t.InfoHash()).String()
}

func Files(t *torrent.Torrent) []*torrent.File {
	if GotInfo(t) != nil {
		return nil
	}
	files := t.Files()
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path() < files[j].Path()
	})
	return files
}

func State(t *torrent.Torrent) torrent.TorrentStats {
	return t.Stats()
}
