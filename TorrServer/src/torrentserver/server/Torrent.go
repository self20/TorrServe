package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"torrentserver/db"
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
	Link string `json:",omitempty"`
	Hash string `json:",omitempty"`
}

type TorrentJsonResponse struct {
	Name    string    `json:",omitempty"`
	Magnet  string    `json:",omitempty"`
	Hash    string    `json:",omitempty"`
	Length  int64     `json:",omitempty"`
	AddTime int64     `json:",omitempty"`
	Size    int64     `json:",omitempty"`
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

	if jreq.Link == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Link must be non-empty")
	}

	torr, err := torrent.Add(jreq.Link)
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

	tor, err := torrent.Get(jreq.Hash)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
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
	list, _ := torrent.List()
	for _, tor := range list {
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

	stat, err := torrent.State(jreq.Hash)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	pstat := torrent.GetPreloadStat(jreq.Hash)

	type jsret struct {
		*torrent.TorrentStat
		*torrent.PreloadStat
	}

	ret := new(jsret)
	ret.TorrentStat = stat
	ret.PreloadStat = pstat

	return c.JSON(http.StatusOK, ret)
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

func getTorrentJS(tor *db.Torrent) (*TorrentJsonResponse, error) {
	js := new(TorrentJsonResponse)
	js.Name = tor.Name
	js.Magnet = tor.Magnet
	js.Hash = tor.Hash
	js.AddTime = tor.Timestamp
	js.Size = tor.Size
	var size int64 = 0
	for _, f := range tor.Files {
		size += f.Size
		tf := TorFile{
			Name:   f.Name,
			Link:   "/torrent/view/" + js.Hash + "/" + utils.FileToLink(f.Name),
			Size:   f.Size,
			Viewed: f.Viewed,
		}
		js.Files = append(js.Files, tf)
	}
	js.Length = size
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
