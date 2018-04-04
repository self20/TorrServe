package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	gotorrent "github.com/anacrolix/torrent"
	"torrentserver/settings"
	"torrentserver/torrent"
	"torrentserver/utils"

	"github.com/labstack/echo"
)

func initTorrent(e *echo.Echo) {

	torrent.Connect()

	e.POST("/torrent/add", torrentAdd)
	e.POST("/torrent/get", torrentGet)
	e.POST("/torrent/rem", torrentRem)
	e.POST("/torrent/list", torrentList)
	e.POST("/torrent/stat", torrentStat)

	e.POST("/torrent/cleancache", torrentCleanCache)

	e.GET("/torrent/view/:hash/:file", torrentView)
	e.HEAD("/torrent/view/:hash/:file", torrentView)
}

type TorrentJsonRequest struct {
	Magnet string `json:",omitempty"`
	Hash   string `json:",omitempty"`
}

type TorrentJsonResponse struct {
	Name    string    `json:",omitempty"`
	Magnet  string    `json:",omitempty"`
	Hash    string    `json:",omitempty"`
	Length  int64     `json:",omitempty"`
	AddTime int64     `json:",omitempty"`
	Files   []TorFile `json:",omitempty"`
}

type TorFile struct {
	Name   string
	Link   string
	Size   int64
	Viewed bool
}

func torrentAdd(c echo.Context) error {
	jreq, err := getJsReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if jreq.Magnet == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Magnet must be non-empty")
	}

	torr, err := torrent.Add(jreq.Magnet)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	js, err := getTorrentJS(torr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, js)
}

func torrentGet(c echo.Context) error {
	jreq, err := getJsReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if jreq.Hash == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Hash must be non-empty")
	}

	tor := torrent.Get(jreq.Hash)
	js, err := getTorrentJS(tor)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, js)
}

func torrentRem(c echo.Context) error {
	jreq, err := getJsReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if jreq.Hash == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Hash must be non-empty")
	}

	torrent.Rem(jreq.Hash)

	return c.JSON(http.StatusOK, nil)
}

func torrentList(c echo.Context) error {
	js := make([]TorrentJsonResponse, 0)
	for _, tor := range torrent.List() {
		jsTor, err := getTorrentJS(tor)
		if err != nil {
			fmt.Println("Error get torrent:", err)
		} else {
			js = append(js, *jsTor)
		}
	}
	sort.Slice(js, func(i, j int) bool {
		if js[i].AddTime == js[j].AddTime {
			return js[i].Name < js[j].Name
		}
		return js[i].AddTime > js[j].AddTime
	})
	return c.JSON(http.StatusOK, js)
}

func torrentStat(c echo.Context) error {
	jreq, err := getJsReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if jreq.Hash == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Hash must be non-empty")
	}

	tor := torrent.Get(jreq.Hash)
	return c.JSON(http.StatusOK, tor.Stats())
}

func torrentCleanCache(c echo.Context) error {
	if err := torrent.Connect(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Torrent server not started: "+err.Error())
	}

	torrent.CleanCache()
	return c.JSON(http.StatusOK, nil)
}

func torrentView(c echo.Context) error {
	hash, err := url.PathUnescape(c.Param("hash"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fileLink, err := url.PathUnescape(c.Param("file"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return torrent.Play(hash, fileLink, c)
}

func getTorrentJS(tor *gotorrent.Torrent) (*TorrentJsonResponse, error) {
	js := new(TorrentJsonResponse)
	js.Name = tor.Name()
	js.Magnet = torrent.Magnet(tor)
	js.Hash = tor.InfoHash().HexString()
	js.AddTime, _ = torrent.GetTime(tor)
	var size int64 = 0
	if torrent.GotInfo(tor) == nil {
		for _, f := range torrent.Files(tor) {
			size += f.Length()
			viewed, _ := settings.ExistTorrView(tor.InfoHash().HexString(), f.Path())
			tf := TorFile{
				Name:   f.Path(),
				Link:   "/torrent/view/" + js.Hash + "/" + utils.FileToLink(f.Path()),
				Size:   f.Length(),
				Viewed: viewed,
			}
			js.Files = append(js.Files, tf)
		}
		js.Length = size
	}
	return js, nil
}

func getJsReq(c echo.Context) (*TorrentJsonRequest, error) {
	decoder := json.NewDecoder(c.Request().Body)
	js := new(TorrentJsonRequest)
	err := decoder.Decode(js)
	if err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, offset=%v", ute.Type, ute.Value, ute.Offset))
		} else if se, ok := err.(*json.SyntaxError); ok {
			return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
		} else {
			return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	return js, nil
}
