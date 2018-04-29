package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"sort"
	"strconv"

	"torrentserver/server/templates"
	"torrentserver/torrent"
	"torrentserver/version"

	"github.com/anacrolix/sync"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/bytes"
)

var (
	server  *echo.Echo
	mutex   sync.Mutex
	fnMutex sync.Mutex
	err     error
)

func Start() {
	if server != nil {
		Stop()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Start web server, version:", version.Version)
	mutex.Lock()
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	server = echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.HTTPErrorHandler = HTTPErrorHandler

	//server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	templates.InitTemplate(server)
	initTorrent(server)
	initSettings(server)

	server.GET("/", mainPage)
	server.GET("/echo", echoPage)
	server.GET("/cache", cachePage)
	server.GET("/stat", statePage)

	go func() {

		server.Listener, err = net.Listen("tcp", "0.0.0.0:8090")
		if err != nil {
			return
		}

		err = server.Start("0.0.0.0:8090")
		server = nil
		mutex.Unlock()
	}()
}

func Stop() {
	fnMutex.Lock()
	defer fnMutex.Unlock()
	if server != nil {
		fmt.Println("Stop web server")
		go torrent.Disconnect()
		server.Close()
		server = nil
	}
}

func Wait() error {
	mutex.Lock()
	mutex.Unlock()
	return err
}

func mainPage(c echo.Context) error {
	return c.Render(http.StatusOK, "mainPage", nil)
}

func echoPage(c echo.Context) error {
	return c.String(http.StatusOK, "Ok")
}

func cachePage(c echo.Context) error {
	infoStates := torrent.CacheState()

	msg := ""

	for _, info := range infoStates {
		if info.Filled == 0 {
			continue
		}
		msg += fmt.Sprintf("Hash: %v\n", info.Hash)
		msg += fmt.Sprintf("Capacity: %d (%v)\n", info.Capacity, bytes.Format(int64(info.Capacity)))
		msg += fmt.Sprintf("Current Size: %d (%v)\n", info.Filled, bytes.Format(int64(info.Filled)))
		msg += fmt.Sprintf("Piece read: %d - %d of %d\n", info.CurrentRead, info.CurrentRead+(info.Capacity/info.PiecesLength), info.PiecesCount)

		msg += "Pieces for delete\n"
		for _, item := range info.PiecesForDel {
			msg += fmt.Sprintf("Hash: %v \t Access: %s\t Buffer size: %d(%s)\t Complete: %v \t Hash: %s\n", item.Id, item.Accessed.Format("15:04:05.000"), item.BufferSize, bytes.Format(int64(item.BufferSize)), item.Completed, item.Hash)
		}

		msg += "\nPieces with buffer\n"
		for _, item := range info.PiecesInCache {
			msg += fmt.Sprintf("Hash: %v \t Access: %s\t Buffer size: %d(%s)\t Complete: %v \t Hash: %s\n", item.Id, item.Accessed.Format("15:04:05.000"), item.BufferSize, bytes.Format(int64(item.BufferSize)), item.Completed, item.Hash)
		}

		msg += "\n"
	}
	return c.String(http.StatusOK, msg)
}

func statePage(c echo.Context) error {
	torrs, err := torrent.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	msg := ""

	sort.Slice(torrs, func(i, j int) bool {
		return torrs[i].Name < torrs[j].Name
	})

	for _, tor := range torrs {
		st, err := torrent.State(tor.Hash)
		if err == nil {
			msg += fmt.Sprintf("Torrent: %v<br>\n", tor.Name)
			msg += fmt.Sprintf("Hash: %v<br>\n<br>\n", tor.Hash)

			msg += fmt.Sprintf("TotalPeers:   	 %v<br>\n", st.TotalPeers)
			msg += fmt.Sprintf("PendingPeers: 	 %v<br>\n", st.PendingPeers)
			msg += fmt.Sprintf("ActivePeers:      %v<br>\n", st.ActivePeers)
			msg += fmt.Sprintf("ConnectedSeeders: %v<br>\n", st.ConnectedSeeders)
			msg += fmt.Sprintf("HalfOpenPeers: 	 %v<br>\n<br>\n", st.HalfOpenPeers)

			msg += fmt.Sprintf("BytesWritten:     %v<br>\n", bytes.Format(st.BytesWritten))
			msg += fmt.Sprintf("BytesWrittenData: %v<br>\n<br>\n", bytes.Format(st.BytesWrittenData))

			msg += fmt.Sprintf("BytesRead: 			%v<br>\n", bytes.Format(st.BytesRead))
			msg += fmt.Sprintf("BytesReadData: 		%v<br>\n", bytes.Format(st.BytesReadData))
			msg += fmt.Sprintf("BytesReadUsefulData: %v<br>\n<br>\n", bytes.Format(st.BytesReadUsefulData))

			msg += fmt.Sprintf("Download Speed: %v/Sec<br>\n", Format(st.DownloadSpeed))
			msg += fmt.Sprintf("Upload Speed:   %v/Sec<br>\n<br>\n", Format(st.UploadSpeed))

			msg += fmt.Sprintf("ChunksWritten:      %v<br>\n", st.ChunksWritten)
			msg += fmt.Sprintf("ChunksRead: 	       %v<br>\n", st.ChunksRead)
			msg += fmt.Sprintf("ChunksReadUseful:   %v<br>\n", st.ChunksReadUseful)
			msg += fmt.Sprintf("ChunksReadUnwanted: %v<br>\n<br>\n", st.ChunksReadUnwanted)

			msg += fmt.Sprintf("PiecesDirtiedGood: %v<br>\n", st.PiecesDirtiedGood)
			msg += fmt.Sprintf("PiecesDirtiedBad:  %v<br>\n", st.PiecesDirtiedBad)

			pstat := torrent.GetPreloadStat(tor.Hash)
			if pstat != nil && pstat.IsPreload {
				msg += fmt.Sprintf("<br>\nPreloadFile: %v<br>\n", pstat.PreloadFile)
				msg += fmt.Sprintf("PreloadOffset:  %v<br>\n", bytes.Format(pstat.PreloadOffset))
				msg += fmt.Sprintf("PreloadLength:  %v<br>\n", bytes.Format(pstat.PreloadLength))
			}

			msg += "<hr><br>\n\n"
		}
	}
	if msg == "" {
		msg = "No connected torrents"
	}
	return c.HTML(http.StatusOK, msg)
}

func HTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			msg = fmt.Sprintf("%v, %v", err, he.Internal)
		}
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = echo.Map{"message": msg}
	}

	log.Println("Web server error:", err, c.Request().URL)

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

const (
	_ = 1.0 << (10 * iota) // ignore first value by assigning to blank identifier
	KB
	MB
	GB
	TB
	PB
	EB
)

func Format(b float64) string {
	multiple := ""
	value := b

	switch {
	case b >= EB:
		value /= EB
		multiple = "EB"
	case b >= PB:
		value /= PB
		multiple = "PB"
	case b >= TB:
		value /= TB
		multiple = "TB"
	case b >= GB:
		value /= GB
		multiple = "GB"
	case b >= MB:
		value /= MB
		multiple = "MB"
	case b >= KB:
		value /= KB
		multiple = "KB"
	case b == 0:
		return "0"
	default:
		return strconv.FormatInt(int64(b), 10) + "B"
	}

	return fmt.Sprintf("%.2f%s", value, multiple)
}
