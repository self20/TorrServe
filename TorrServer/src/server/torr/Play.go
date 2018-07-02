package torr

import (
	"fmt"
	"net/http"
	"time"

	"server/settings"

	"github.com/anacrolix/missinggo/httptoo"
	"github.com/anacrolix/torrent"
	"github.com/labstack/echo"
)

func (bt *BTServer) Play(state *TorrentState, file *torrent.File, timestamp time.Time, c echo.Context) error {
	bt.watcher()
	go settings.SetViewed(state.Hash, file.Path())
	reader := file.NewReader()
	reader.SetReadahead(getReadahead())

	state.readers++

	fmt.Println("Connect reader:", state.readers)
	c.Response().Header().Set("Connection", "close")
	c.Response().Header().Set("ETag", httptoo.EncodeQuotedString(fmt.Sprintf("%s/%s", state.Hash, file.Path())))

	defer func() {
		fmt.Println("Disconnect reader:", state.readers)
		reader.Close()
		state.expiredTime = time.Now().Add(time.Second * 20)
		state.readers--
		go bt.watcher()
	}()

	http.ServeContent(c.Response(), c.Request(), file.Path(), timestamp, reader)
	return c.NoContent(http.StatusOK)
}

func getReadahead() int64 {
	readahead := int64(float64(settings.Get().CacheSize) * 0.33)
	if readahead < 66*1024*1024 {
		readahead = int64(settings.Get().CacheSize)
		if readahead > 66*1024*1024 {
			readahead = 66 * 1024 * 1024
		}
	}
	return readahead
}
