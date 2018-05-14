package torr

import (
	"fmt"
	"net/http"
	"time"

	"server/settings"
	"server/utils"

	"github.com/anacrolix/torrent"
	"github.com/labstack/echo"
)

func (bt *BTServer) Play(hash, fileLink string, c echo.Context) error {
	tordb, err := bt.Get(hash)
	if err != nil {
		return c.String(http.StatusNotFound, "Torrent not found:"+err.Error()+" "+hash+"/"+fileLink)
	}

	var file *settings.File
	for _, f := range tordb.Files {
		if utils.FileToLink(f.Name) == fileLink {
			file = &f
			break
		}
	}

	if file == nil {
		return c.String(http.StatusNotFound, "File in torrent not found: "+hash+"/"+fileLink)
	}

	state, err := bt.getTorrent(tordb)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	go settings.SetViewed(tordb.Hash, file.Name)

	var reader torrent.Reader
	for _, f := range state.torrent.Files() {
		if f.Path() == file.Name {
			reader = f.NewReader()
			break
		}
	}

	readahead := int64(float64(settings.Get().CacheSize) * 0.33)
	if readahead < 5*1024*1024 {
		readahead = 5 * 1024 * 1024
	}
	reader.SetReadahead(readahead)
	state.readers++
	fmt.Println("Connect reader:", state.readers)

	tm := settings.StartTime
	if tordb.Timestamp != 0 {
		tm = time.Unix(tordb.Timestamp, 0)
	}

	utils.ServeContentTorrent(c.Response(), c.Request(), file.Name, tm, file.Size, reader)
	reader.Close()

	fmt.Println("Disconnect reader:", state.readers)
	state.readers--

	if state.readers == 0 {
		state.expiredTime = time.Now().Add(time.Minute)
	}

	return c.NoContent(http.StatusOK)
}
