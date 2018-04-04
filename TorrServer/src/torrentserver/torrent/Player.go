package torrent

import (
	"net/http"
	"time"

	"torrentserver/settings"
	"torrentserver/utils"

	"github.com/anacrolix/torrent"
	"github.com/labstack/echo"
)

var (
	currentTorrent *torrent.Torrent
)

func Play(hash, fileLink string, c echo.Context) error {
	tor := Get(hash)
	if tor == nil {
		return c.String(http.StatusNotFound, "Torrent not found: "+hash+"/"+fileLink)
	}

	if currentTorrent == nil || tor.InfoHash() != currentTorrent.InfoHash() {
		if currentTorrent != nil {
			go Stop(currentTorrent)
		}
		currentTorrent = tor

		if err := GotInfo(tor); err != nil {
			return err
		}
	}

	var file *torrent.File
	for _, f := range tor.Files() {
		if utils.FileToLink(f.Path()) == fileLink {
			file = f
			break
		}
	}

	if file == nil {
		return c.String(http.StatusNotFound, "File in torrent not found: "+hash+"/"+fileLink)
	}

	go settings.SaveTorrentView(tor.InfoHash().HexString(), file.Path())

	reader := NewReader(tor, file)

	tmi, _ := GetTime(tor)
	tm := settings.StartTime
	if tmi != 0 {
		tm = time.Unix(tmi, 0)
	}

	utils.ServeContentTorrent(c.Response(), c.Request(), tor.Name(), tm, file.Length(), reader)

	reader.Close()
	return c.JSON(http.StatusOK, nil)
}

func Stop(tor *torrent.Torrent) {
	if tor == nil {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	info := tor.Metainfo()
	tor.Drop()
	client.AddTorrent(&info)
}
