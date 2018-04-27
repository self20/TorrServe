package torrent

import (
	"net/http"
	"time"

	"torrentserver/db"
	"torrentserver/settings"
	"torrentserver/utils"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/labstack/echo"
)

var (
	currentTorrent *db.Torrent
)

func Play(hash, fileLink string, c echo.Context) error {
	tordb, err := Get(hash)
	if err != nil {
		return c.String(http.StatusNotFound, "Torrent not found:"+err.Error()+" "+hash+"/"+fileLink)
	}

	var file *db.File
	for _, f := range tordb.Files {
		if utils.FileToLink(f.Name) == fileLink {
			file = &f
			break
		}
	}

	if file == nil {
		return c.String(http.StatusNotFound, "PreloadFile in torrent not found: "+hash+"/"+fileLink)
	}

	go db.SetViewed(tordb.Hash, file.Name)

	reader, err := handler.NewReader(tordb, file.Name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	mhash := metainfo.NewHashFromHex(tordb.Hash)
	hl := handler.GetHandle(mhash)
	if hl != nil && hl.Preload != nil && hl.Preload.preload {
		for hl.Preload.preload {
			time.Sleep(time.Millisecond * 100)
		}
	}

	tm := settings.StartTime
	if tordb.Timestamp != 0 {
		tm = time.Unix(tordb.Timestamp, 0)
	}

	utils.ServeContentTorrent(c.Response(), c.Request(), tordb.Name, tm, file.Size, reader)

	reader.Close()
	return c.JSON(http.StatusOK, nil)
}
