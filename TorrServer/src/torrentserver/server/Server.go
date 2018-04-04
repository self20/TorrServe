package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"sort"

	"torrentserver/server/templates"
	"torrentserver/settings"
	"torrentserver/torrent"
	"torrentserver/utils"

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
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Start web server, version:", utils.Version)
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
		go torrent.Disconnect()
		server.Close()
		server = nil
		settings.CloseDB()
	}
}

func Wait() error {
	mutex.Lock()
	mutex.Unlock()
	return err
}

func mainPage(c echo.Context) error {
	return c.Render(http.StatusOK, "mainPage", torrent.List())
}

func echoPage(c echo.Context) error {
	return c.String(http.StatusOK, "Ok")
}

func cachePage(c echo.Context) error {
	infoStates := torrent.GetStates()

	msg := ""

	for _, info := range infoStates {
		if info.Filled == 0 {
			continue
		}
		msg += fmt.Sprintf("Hash: %v\n", info.Hash)
		msg += fmt.Sprintf("Capacity: %d (%v)\n", info.Capacity, bytes.Format(int64(info.Capacity)))
		msg += fmt.Sprintf("Current Size: %d (%v)\n", info.Filled, bytes.Format(int64(info.Filled)))
		msg += fmt.Sprintf("Piece read: %d - %d of %d\n", info.CurrentRead, info.CurrentRead+(info.Capacity/info.PiecesLength), info.PiecesCount)

		for _, item := range info.Pieces {
			msg += fmt.Sprintf("Hash: %v \t Access: %s\t Buffer size: %d(%s)\t Complete: %v \t Hash: %s\n", item.Id, item.Accessed.Format("15:04:05.000"), item.BufferSize, bytes.Format(int64(item.BufferSize)), item.Completed, item.Hash)
		}
		msg += "\n"
	}
	return c.String(http.StatusOK, msg)
}

func statePage(c echo.Context) error {
	torrs := torrent.List()

	msg := ""

	sort.Slice(torrs, func(i, j int) bool {
		return torrs[i].Name() < torrs[j].Name()
	})

	for _, tor := range torrs {
		st := tor.Stats()

		msg += fmt.Sprintf("Torrent: %v\n", tor.Name())
		msg += fmt.Sprintf("TotalPeers: %v\n", st.TotalPeers)
		msg += fmt.Sprintf("PendingPeers: %v\n", st.PendingPeers)
		msg += fmt.Sprintf("ActivePeers: %v\n", st.ActivePeers)
		msg += fmt.Sprintf("ConnectedSeeders: %v\n", st.ConnectedSeeders)
		msg += fmt.Sprintf("HalfOpenPeers: %v\n", st.HalfOpenPeers)

		msg += fmt.Sprintf("BytesWritten: %v\n", st.BytesWritten)
		msg += fmt.Sprintf("BytesWrittenData: %v\n", st.BytesWrittenData)

		msg += fmt.Sprintf("BytesRead: %v\n", st.BytesRead)
		msg += fmt.Sprintf("BytesReadData: %v\n", st.BytesReadData)
		msg += fmt.Sprintf("BytesReadUsefulData: %v\n", st.BytesReadUsefulData)

		msg += fmt.Sprintf("ChunksWritten: %v\n", st.ChunksWritten)

		msg += fmt.Sprintf("ChunksRead: %v\n", st.ChunksRead)
		msg += fmt.Sprintf("ChunksReadUseful: %v\n", st.ChunksReadUseful)
		msg += fmt.Sprintf("ChunksReadUnwanted: %v\n", st.ChunksReadUnwanted)

		msg += fmt.Sprintf("PiecesDirtiedGood: %v\n", st.PiecesDirtiedGood)
		msg += fmt.Sprintf("PiecesDirtiedBad: %v\n", st.PiecesDirtiedBad)

		msg += "\n"
	}
	return c.String(http.StatusOK, msg)
}

func HTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Inner != nil {
			msg = fmt.Sprintf("%v, %v", err, he.Inner)
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
