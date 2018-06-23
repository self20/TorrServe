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

func (bt *BTServer) Play(state *TorrentState, file *torrent.File, timestamp time.Time, c echo.Context) error {
	bt.watcher()
	go settings.SetViewed(state.Hash, file.Path())
	reader := file.NewReader()
	//readahead := int64(float64(settings.Get().CacheSize) * 0.33)
	//if readahead < 5*1024*1024 {
	//	readahead = 5 * 1024 * 1024
	//}

	readahead := int64(float64(settings.Get().CacheSize) * 0.33)
	if readahead < 66*1024*1024 {
		readahead = int64(settings.Get().CacheSize)
		if readahead > 66*1024*1024 {
			readahead = 66 * 1024 * 1024
		}
	}
	reader.SetReadahead(readahead)

	state.readers++
	defer func() {
		state.expiredTime = time.Now().Add(time.Second * 20)
		state.readers--
		go bt.watcher()
	}()

	fmt.Println("Connect reader:", state.readers)

	utils.ServeContentTorrent(c.Response(), c.Request(), file.Path(), timestamp, file.Length(), reader)
	reader.Close()

	fmt.Println("Disconnect reader:", state.readers)

	return c.NoContent(http.StatusOK)
}
