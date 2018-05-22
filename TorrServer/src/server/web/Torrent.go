package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"time"

	"server/settings"
	"server/utils"

	"github.com/labstack/echo"
)

func initTorrent(e *echo.Echo) {
	e.POST("/torrent/add", torrentAdd)
	e.POST("/torrent/get", torrentGet)
	e.POST("/torrent/rem", torrentRem)
	e.POST("/torrent/list", torrentList)
	e.POST("/torrent/stat", torrentStat)

	e.POST("/torrent/cleancache", torrentCleanCache)
	e.GET("/torrent/restart", torrentRestart)

	e.GET("/torrent/playlist/:hash/*", torrentPlayList)
	e.GET("/torrent/playlist.m3u", torrentPlayListAll)

	e.GET("/torrent/view/:hash/:file", torrentView)
	e.HEAD("/torrent/view/:hash/:file", torrentViewHead)
	e.GET("/torrent/preload/:hash/:file", torrentPreload)
}

type TorrentJsonRequest struct {
	Link string `json:",omitempty"`
	Hash string `json:",omitempty"`
}

type TorrentJsonResponse struct {
	Name     string    `json:",omitempty"`
	Magnet   string    `json:",omitempty"`
	Hash     string    `json:",omitempty"`
	Length   int64     `json:",omitempty"`
	AddTime  int64     `json:",omitempty"`
	Size     int64     `json:",omitempty"`
	Playlist string    `json:",omitempty"`
	Files    []TorFile `json:",omitempty"`
}

type TorFile struct {
	Name     string
	Link     string
	Playlist string
	Size     int64
	Viewed   bool
}

func torrentAdd(c echo.Context) error {
	jreq, err := getJsReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if jreq.Link == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Link must be non-empty")
	}

	torr, err := bts.Add(jreq.Link)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	js, err := getTorrentJS(torr)
	if err != nil {
		fmt.Println("Error add torrent:", torr.Hash, err)
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

	torr, err := bts.Get(jreq.Hash)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	js, err := getTorrentJS(torr)
	if err != nil {
		fmt.Println("Error add torrent:", torr.Hash, err)
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

	bts.Rem(jreq.Hash)

	return c.JSON(http.StatusOK, nil)
}

func torrentList(c echo.Context) error {
	js := make([]TorrentJsonResponse, 0)
	list, _ := bts.List()
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

	stat := bts.TorrentState(jreq.Hash)
	if stat == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, stat)
}

func torrentPreload(c echo.Context) error {
	if settings.Get().PreloadBufferSize > 0 {

		hash, err := url.PathUnescape(c.Param("hash"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fileLink, err := url.PathUnescape(c.Param("file"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if hash == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Hash must be non-empty")
		}
		if fileLink == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "File link must be non-empty")
		}

		err = bts.Preload(hash, fileLink)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}

func torrentCleanCache(c echo.Context) error {
	jreq, err := getJsReq(c)
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

	torr, err := bts.Get(hash)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	m3u := "#EXTM3U\n"
	m3u += "#EXT-X-PLAYLIST-TYPE:VOD\n"
	m3u += "#EXT-X-VERSION:3\n"

	for _, f := range torr.Files {
		m3u += "#EXTINF:-1," + f.Name + "\n"
		m3u += c.Scheme() + "://" + c.Request().Host + "/torrent/view/" + hash + "/" + utils.FileToLink(f.Name) + "\n\n"
	}

	http.ServeContent(c.Response(), c.Request(), torr.Name+".m3u", time.Now(), bytes.NewReader([]byte(m3u)))
	return c.NoContent(http.StatusOK)
}

func torrentPlayListAll(c echo.Context) error {
	list, err := bts.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	m3u := "#EXTM3U\n"
	m3u += "#EXT-X-PLAYLIST-TYPE:VOD\n"
	m3u += "#EXT-X-VERSION:3\n"

	for _, t := range list {
		m3u += "#EXTINF:0," + t.Name + "\n"
		m3u += c.Scheme() + "://" + c.Request().Host + "/torrent/playlist/" + t.Hash + "/" + utils.FileToLink(t.Name) + ".m3u" + "\n\n"
	}

	http.ServeContent(c.Response(), c.Request(), "playlist.m3u", time.Now(), bytes.NewReader([]byte(m3u)))
	return c.NoContent(http.StatusOK)
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

	return bts.Play(hash, fileLink, c)
}

func torrentViewHead(c echo.Context) error {
	hash, err := url.PathUnescape(c.Param("hash"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fileLink, err := url.PathUnescape(c.Param("file"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tordb, err := bts.Get(hash)
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

	tm := settings.StartTime
	if tordb.Timestamp != 0 {
		tm = time.Unix(tordb.Timestamp, 0)
	}

	ctype := mime.TypeByExtension(filepath.Ext(file.Name))
	if ctype == "" {
		ctype = utils.GetMimeType(file.Name)
	}

	c.Response().Header().Set("Accept-Ranges", "bytes")
	c.Response().Header().Set("Content-Length", fmt.Sprint(file.Size))
	c.Response().Header().Set("Content-Type", ctype)
	c.Response().Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	c.Response().Header().Set("Last-Modified", tm.UTC().Format(http.TimeFormat))

	return c.NoContent(http.StatusOK)
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

func getJsReq(c echo.Context) (*TorrentJsonRequest, error) {
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
