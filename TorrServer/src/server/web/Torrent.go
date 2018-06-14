package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"

	"server/settings"
	"server/torr"
	"server/utils"
	"server/web/helpers"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/labstack/echo"
)

func initTorrent(e *echo.Echo) {
	e.POST("/torrent/add", torrentAdd)
	e.POST("/torrent/upload", torrentUpload)
	e.POST("/torrent/get", torrentGet)
	e.POST("/torrent/rem", torrentRem)
	e.POST("/torrent/list", torrentList)
	e.POST("/torrent/stat", torrentStat)

	e.POST("/torrent/cleancache", torrentCleanCache)
	e.GET("/torrent/restart", torrentRestart)

	e.GET("/torrent/playlist/:hash/*", torrentPlayList)
	e.GET("/torrent/playlist.m3u", torrentPlayListAll)

	e.GET("/torrent/view/:hash/:file", torrentView)
	e.HEAD("/torrent/view/:hash/:file", torrentView)
	e.GET("/torrent/preload/:hash/:file", torrentPreload)
}

type TorrentJsonRequest struct {
	Link     string `json:",omitempty"`
	Hash     string `json:",omitempty"`
	DontSave bool   `json:",omitempty"`
}

type TorrentJsonResponse struct {
	Name          string
	Magnet        string
	Hash          string
	Length        int64
	AddTime       int64
	Size          int64
	IsGettingInfo bool
	Playlist      string
	Files         []TorFile `json:",omitempty"`
}

type TorFile struct {
	Name     string
	Link     string
	Playlist string
	Size     int64
	Viewed   bool
}

func torrentAdd(c echo.Context) error {
	jreq, err := getJsReqTorr(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if jreq.Link == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Link must be non-empty")
	}

	magnet, err := helpers.GetMagnet(jreq.Link)
	if err != nil {
		fmt.Println("Error get magnet:", jreq.Hash, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = helpers.Add(bts, magnet, !jreq.DontSave)
	if err != nil {
		fmt.Println("Error add torrent:", jreq.Hash, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, magnet.InfoHash.HexString())
}

func torrentUpload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer form.RemoveAll()

	_, dontSave := form.Value["DontSave"]
	var magnets []metainfo.Magnet

	for _, file := range form.File {
		torrFile, err := file[0].Open()
		if err != nil {
			return err
		}
		defer torrFile.Close()

		mi, err := metainfo.Load(torrFile)
		if err != nil {
			fmt.Println("Error upload torrent", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		info, err := mi.UnmarshalInfo()
		if err != nil {
			fmt.Println("Error upload torrent", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		magnet := mi.Magnet(info.Name, mi.HashInfoBytes())
		magnets = append(magnets, magnet)
	}

	ret := make([]string, 0)
	for _, magnet := range magnets {
		er := helpers.Add(bts, &magnet, !dontSave)
		if er != nil {
			err = er
			fmt.Println("Error add torrent:", magnet.String(), er)
		}
		ret = append(ret, magnet.InfoHash.HexString())
	}

	return c.JSON(http.StatusOK, ret)
}

func torrentGet(c echo.Context) error {
	jreq, err := getJsReqTorr(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if jreq.Hash == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Hash must be non-empty")
	}

	torr, err := settings.LoadTorrentDB(jreq.Hash)
	if err != nil {
		fmt.Println("Error get torrent:", jreq.Hash, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	isGotInfo := false
	if torr == nil {
		hash := metainfo.NewHashFromHex(jreq.Hash)
		ts := bts.GetTorrent(hash)
		isGotInfo = ts.IsGettingInfo
		torr = toTorrentDB(ts)
	}

	if torr == nil {
		fmt.Println("Error get: torrent not found", jreq.Hash)
		return echo.NewHTTPError(http.StatusBadRequest, "Error get: torrent not found "+jreq.Hash)
	}

	js, err := getTorrentJS(torr)
	if err != nil {
		fmt.Println("Error get torrent:", torr.Hash, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	js.IsGettingInfo = isGotInfo
	return c.JSON(http.StatusOK, js)
}

func torrentRem(c echo.Context) error {
	jreq, err := getJsReqTorr(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if jreq.Hash == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Hash must be non-empty")
	}

	settings.RemoveTorrentDB(jreq.Hash)
	bts.RemoveTorrent(jreq.Hash)

	return c.JSON(http.StatusOK, nil)
}

func torrentList(c echo.Context) error {
	js := make([]TorrentJsonResponse, 0)
	list, _ := settings.LoadTorrentsDB()

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

	slist := bts.List()

	find := func(tsr []TorrentJsonResponse, st *torr.TorrentState) bool {
		for _, ts := range tsr {
			if ts.Hash == st.Hash {
				return true
			}
		}
		return false
	}

	for _, st := range slist {
		if !find(js, st) {
			tdb := toTorrentDB(st)
			jsTor, err := getTorrentJS(tdb)
			if err != nil {
				fmt.Println("Error get torrent:", err)
			} else {
				jsTor.IsGettingInfo = st.IsGettingInfo
				js = append(js, *jsTor)
			}
		}
	}

	return c.JSON(http.StatusOK, js)
}

func torrentStat(c echo.Context) error {
	jreq, err := getJsReqTorr(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if jreq.Hash == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Hash must be non-empty")
	}

	hash := metainfo.NewHashFromHex(jreq.Hash)
	stat := bts.GetTorrent(hash)
	if stat == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, stat)
}

func torrentPreload(c echo.Context) error {
	if settings.Get().PreloadBufferSize > 0 {
		hashHex, err := url.PathUnescape(c.Param("hash"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fileLink, err := url.PathUnescape(c.Param("file"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if hashHex == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Hash must be non-empty")
		}
		if fileLink == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "File link must be non-empty")
		}

		hash := metainfo.NewHashFromHex(hashHex)
		st := bts.GetTorrent(hash)
		if st == nil {
			torrDb, err := settings.LoadTorrentDB(hashHex)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Torrent not found: "+hashHex)
			}
			m, err := metainfo.ParseMagnetURI(torrDb.Magnet)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Error parser magnet in db: "+hashHex)
			}
			st, err = bts.AddTorrent(&m)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		}
		file := helpers.FindFile(fileLink, st.Torrent)

		err = bts.Preload(hash, file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}

func torrentCleanCache(c echo.Context) error {
	jreq, err := getJsReqTorr(c)
	if err != nil {
		bts.Clean("")
		return c.JSON(http.StatusOK, nil)
	}
	bts.Clean(jreq.Hash)
	return c.JSON(http.StatusOK, nil)
}

func torrentRestart(c echo.Context) error {
	fmt.Println("Restart torrent engine")
	err := bts.Reconnect()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Ok")
}

func torrentPlayList(c echo.Context) error {
	hash, err := url.PathUnescape(c.Param("hash"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	torr, err := settings.LoadTorrentDB(hash)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	m3u := "#EXTM3U\n"

	for _, f := range torr.Files {
		if utils.GetMimeType(f.Name) != "*/*" {
			m3u += "#EXTINF:-1," + f.Name + "\n"
			m3u += c.Scheme() + "://" + c.Request().Host + "/torrent/view/" + hash + "/" + utils.FileToLink(f.Name) + "\n\n"
		}
	}

	c.Response().Header().Set("Content-Type", "audio/x-mpegurl")
	http.ServeContent(c.Response(), c.Request(), torr.Name+".m3u", time.Now(), bytes.NewReader([]byte(m3u)))
	return c.NoContent(http.StatusOK)
}

func torrentPlayListAll(c echo.Context) error {
	list, err := settings.LoadTorrentsDB()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	m3u := "#EXTM3U\n"

	for _, t := range list {
		m3u += "#EXTINF:0," + t.Name + "\n"
		m3u += c.Scheme() + "://" + c.Request().Host + "/torrent/playlist/" + t.Hash + "/" + utils.FileToLink(t.Name) + ".m3u" + "\n\n"
	}

	c.Response().Header().Set("Content-Type", "audio/x-mpegurl")
	http.ServeContent(c.Response(), c.Request(), "playlist.m3u", time.Now(), bytes.NewReader([]byte(m3u)))
	return c.NoContent(http.StatusOK)
}

func torrentView(c echo.Context) error {
	hashHex, err := url.PathUnescape(c.Param("hash"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fileLink, err := url.PathUnescape(c.Param("file"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	timestamp := settings.StartTime
	hash := metainfo.NewHashFromHex(hashHex)
	st := bts.GetTorrent(hash)
	if st == nil {
		torrDb, err := settings.LoadTorrentDB(hashHex)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Torrent not found: "+hashHex)
		}

		if torrDb.Timestamp != 0 {
			timestamp = time.Unix(torrDb.Timestamp, 0)
		}

		m, err := metainfo.ParseMagnetURI(torrDb.Magnet)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error parser magnet in db: "+hashHex)
		}

		st, err = bts.AddTorrent(&m)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	file := helpers.FindFile(fileLink, st.Torrent)

	return bts.Play(st, file, timestamp, c)
}

func toTorrentDB(ts *torr.TorrentState) *settings.Torrent {
	if ts == nil {
		return nil
	}
	tor := new(settings.Torrent)
	tor.Name = ts.Name
	tor.Hash = ts.Hash
	tor.Timestamp = settings.StartTime.Unix()
	mi := ts.Torrent.Metainfo()
	tor.Magnet = mi.Magnet(ts.Name, ts.Torrent.InfoHash()).String()
	tor.Size = ts.TorrentSize
	for _, f := range ts.Files() {
		tf := settings.File{
			Name:   f.Path(),
			Size:   f.Length(),
			Viewed: false,
		}
		tor.Files = append(tor.Files, tf)
	}
	return tor
}

func getTorrentJS(tor *settings.Torrent) (*TorrentJsonResponse, error) {
	js := new(TorrentJsonResponse)
	js.Name = tor.Name
	js.Magnet = tor.Magnet
	js.Hash = tor.Hash
	js.AddTime = tor.Timestamp
	js.Size = tor.Size
	js.Playlist = "/torrent/playlist/" + tor.Hash + "/" + utils.FileToLink(tor.Name) + ".m3u"
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

func getJsReqTorr(c echo.Context) (*TorrentJsonRequest, error) {
	buf, _ := ioutil.ReadAll(c.Request().Body)
	jsstr := string(buf)
	decoder := json.NewDecoder(bytes.NewBufferString(jsstr))
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
