package server

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"server/search/parser"
	"server/search/tmdb"
	"server/search/torrent"

	"github.com/labstack/echo"
)

func initSearch(e *echo.Echo) {
	e.GET("/search", searchPage)

	e.GET("/search/movie", searchMovie)
	e.GET("/search/show", searchShow)

	e.GET("/search/movie/:id", getMovie)
	e.GET("/search/show/:id", getShow)

	e.GET("/search/config", searchConfig)
	e.GET("/search/torrent", searchTorrent)
}

func searchPage(c echo.Context) error {
	vt := c.QueryParam("vt")

	if c.QueryParam("language") == "" {
		c.QueryParams().Set("language", "ru")
	}

	var pinfo *PageInfo
	if vt == "show" {
		shows, all := getShows(c)
		pinfo = tvToPageInfo(c, shows, all)
	} else if vt == "torrent" {
		pinfo = new(PageInfo)
		pinfo.IsTorrent = true
		torrs, err := getTorrent(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		for _, t := range torrs {
			ii := new(ItemInfo)
			ii.Name = t.Name
			ii.OriginalName = t.Magnet
			ii.Year = t.Size
			ii.Seasons = t.PeersUl
			ii.Episodes = t.PeersDl
			pinfo.Items = append(pinfo.Items, ii)
		}
		return c.Render(http.StatusOK, "searchPage", pinfo)
	} else {
		movies, all := getMovies(c)
		pinfo = movToPageInfo(c, movies, all)
	}
	for i := time.Now().Year(); i > 1900; i-- {
		pinfo.Years = append(pinfo.Years, i)
	}

	return c.Render(http.StatusOK, "searchPage", pinfo)
}

func searchMovie(c echo.Context) error {
	m, _ := getMovies(c)
	return c.JSON(http.StatusOK, m)
}

func searchShow(c echo.Context) error {
	s, _ := getShows(c)
	return c.JSON(http.StatusOK, s)
}

func getMovie(c echo.Context) error {
	ids := c.Param("id")
	if ids == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "empty id")
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ln := c.QueryParam("lanuage")
	if ln == "" {
		ln = "ru"
	}

	mov := tmdb.GetMovie(id, ln)
	if mov == nil {
		return echo.NewHTTPError(http.StatusNotFound, "ids")
	}
	fixMovies(tmdb.Movies{mov})
	return c.JSON(http.StatusOK, mov)
}

func getShow(c echo.Context) error {
	ids := c.Param("id")
	if ids == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "empty id")
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ln := c.QueryParam("lanuage")
	if ln == "" {
		ln = "ru"
	}

	show := tmdb.GetShow(id, ln)
	if show == nil {
		return echo.NewHTTPError(http.StatusNotFound, "ids")
	}

	fixShows(tmdb.Shows{show})

	return c.JSON(http.StatusOK, show)
}

func searchConfig(c echo.Context) error {
	_, typeReq, _, language, _ := getParams(c)
	switch typeReq {
	case "genres":
		mg := tmdb.GetMovieGenres(language)
		sg := tmdb.GetTVGenres(language)
		return c.JSON(http.StatusOK, struct {
			MovieGenres []*tmdb.Genre
			ShowGenres  []*tmdb.Genre
		}{mg, sg})
	case "config":
		cfg := tmdb.GetConfig()
		cfg.ChangeKeys = nil
		return c.JSON(http.StatusOK, cfg)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "unknown type")
	}
}

func searchTorrent(c echo.Context) error {
	torrs, err := getTorrent(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, torrs)
}

func getTorrent(c echo.Context) ([]*parser.Torrent, error) {
	filter, _ := c.QueryParams()["ft"]
	query := c.QueryParam("query")

	if query == "" {
		return nil, nil
	}
	torrs := torrent.Search(query, filter)
	sort.Slice(torrs, func(i, j int) bool {
		gri := getGrTorr(torrs[i])
		grj := getGrTorr(torrs[j])
		if gri != grj {
			return gri < grj
		}

		if torrs[i].PeersUl == -1 && torrs[j].PeersUl == -1 {
			return torrs[i].Size > torrs[j].Size
		}

		return torrs[i].PeersUl > torrs[j].PeersUl
	})
	return torrs, nil
}

//gr 0 - 50 or more
//gr 1 - 30 - 50 gb
//gr 2 - 15 - 30 gb
//gr 3 - 0 - 15 gb
func getGrTorr(t *parser.Torrent) int {
	szStr := ""
	if t.Size[len(t.Size)-2:] == "GB" {
		szStr = strings.TrimSpace(t.Size[:len(t.Size)-2])
	}
	if t.Size[len(t.Size)-4:] == "ГБ" {
		szStr = strings.TrimSpace(t.Size[:len(t.Size)-4])
	}

	if szStr != "" {
		sz, _ := strconv.ParseFloat(szStr, 32)
		if sz > 50 {
			return 0
		}
		if sz > 30 {
			return 1
		}
		if sz > 15 {
			return 2
		}
		if sz > 5 {
			return 3
		}
		if sz > 1 {
			return 4
		}
		return 5
	}
	return 6
}

func getMovies(c echo.Context) (tmdb.Movies, int) {
	params, typeReq, page, language, query := getParams(c)
	var movies tmdb.Movies
	var all int
	switch typeReq {
	case "discover":
		movies, all = tmdb.DiscoverMovies(params, page)
	case "search":
		movies, all = tmdb.SearchMovies(query, language, page)
	}
	return fixMovies(movies), all
}

func getShows(c echo.Context) (tmdb.Shows, int) {
	params, typeReq, page, language, query := getParams(c)
	var shows tmdb.Shows
	var all int
	switch typeReq {
	case "discover":
		shows, all = tmdb.DiscoverShows(params, page)
	case "search":
		shows, all = tmdb.SearchShows(query, language, page)
	}
	return fixShows(shows), all
}

// ?type=discover&page=1&...
// ?type=search&page=1&query=NAME&...
//return params,type,page,language,query
func getParams(c echo.Context) (map[string]string, string, int, string, string) {
	params := make(map[string]string)
	typeReq := "discover"
	page := 1
	query := ""
	language := "ru"
	for k, v := range c.QueryParams() {
		if k == "language" {
			language = strings.Join(v, ",")
			params[k] = language
		}
		if k == "page" {
			page, _ = strconv.Atoi(v[0])
		} else if k == "type" {
			typeReq = v[0]
		} else if k == "query" {
			query = strings.Join(v, " ")
		} else if k == "vt" {
		} else {
			params[k] = strings.Join(v, ",")
		}
	}
	return params, typeReq, page, language, query
}

func fixMovies(movies tmdb.Movies) tmdb.Movies {
	ret := tmdb.Movies{}
	for _, m := range movies {
		if m == nil {
			continue
		}
		if m.BackdropPath != "" {
			m.BackdropPath = tmdb.ImageURL(m.BackdropPath, "w1280")
		}
		if m.PosterPath != "" {
			m.PosterPath = tmdb.ImageURL(m.PosterPath, "w500")
		}

		for _, i := range m.Images.Backdrops {
			i.FilePath = tmdb.ImageURL(i.FilePath, "w1280")
		}

		for _, i := range m.Images.Posters {
			i.FilePath = tmdb.ImageURL(i.FilePath, "w500")
		}

		for _, i := range m.Images.Stills {
			i.FilePath = tmdb.ImageURL(i.FilePath, "w1280")
		}

		for _, t := range m.Trailers.Youtube {
			t.Source = "https://www.youtube.com/watch?v=" + t.Source
		}
		ret = append(ret, m)
	}
	return ret
}

func fixShows(shows tmdb.Shows) tmdb.Shows {
	ret := tmdb.Shows{}
	for _, s := range shows {
		if s == nil {
			continue
		}
		if s.BackdropPath != "" {
			s.BackdropPath = tmdb.ImageURL(s.BackdropPath, "w1280")
		}
		if s.PosterPath != "" {
			s.PosterPath = tmdb.ImageURL(s.PosterPath, "w500")
		}

		for _, i := range s.Images.Backdrops {
			i.FilePath = tmdb.ImageURL(i.FilePath, "w1280")
		}

		for _, i := range s.Images.Posters {
			i.FilePath = tmdb.ImageURL(i.FilePath, "w500")
		}

		for _, i := range s.Images.Stills {
			i.FilePath = tmdb.ImageURL(i.FilePath, "w1280")
		}
		ret = append(ret, s)
	}

	return ret
}

type PageInfo struct {
	Items []*ItemInfo

	Pages     int
	Genres    []*tmdb.Genre
	Sorts     []string
	Years     []int
	IsTorrent bool
}

type ItemInfo struct {
	ID           int
	Name         string
	OriginalName string

	Overview string
	Genres   []*tmdb.IDName
	Year     string
	Tagline  string

	Poster   string
	Backdrop string
	AllArts  []string

	Seasons  int
	Episodes int
}

func movToPageInfo(c echo.Context, movies tmdb.Movies, all int) *PageInfo {
	pi := new(PageInfo)
	pi.Genres = tmdb.GetMovieGenres(c.QueryParam("language"))
	pi.Sorts = []string{"popularity.asc", "popularity.desc", "release_date.asc", "release_date.desc", "revenue.asc", "revenue.desc", "primary_release_date.asc", "primary_release_date.desc", "original_title.asc", "original_title.desc", "vote_average.asc", "vote_average.desc", "vote_count.asc", "vote_count.desc"}
	for i := time.Now().Year(); i > 1900; i-- {
		pi.Years = append(pi.Years, i)
	}

	limit := tmdb.ResultsPerPage * tmdb.PagesAtOnce
	pi.Pages = all/limit + 1
	if pi.Pages > 1000 {
		pi.Pages = 1000
	}

	for _, m := range movies {
		ii := new(ItemInfo)
		ii.ID = m.ID
		ii.Name = m.Title
		ii.OriginalName = m.OriginalTitle
		ii.Overview = m.Overview
		ii.Genres = m.Genres
		if m.ReleaseDate != "" {
			ii.Year = m.ReleaseDate[:4]
		}
		ii.Tagline = m.TagLine
		ii.Poster = m.PosterPath
		ii.Backdrop = m.BackdropPath
		for _, i := range m.Images.Posters {
			ii.AllArts = append(ii.AllArts, i.FilePath)
		}
		for _, i := range m.Images.Backdrops {
			ii.AllArts = append(ii.AllArts, i.FilePath)
		}
		for _, i := range m.Images.Stills {
			ii.AllArts = append(ii.AllArts, i.FilePath)
		}
		ii.Seasons = 0
		ii.Episodes = 0
		pi.Items = append(pi.Items, ii)
	}
	return pi
}

func tvToPageInfo(c echo.Context, shows tmdb.Shows, all int) *PageInfo {
	pi := new(PageInfo)
	pi.Genres = tmdb.GetTVGenres(c.QueryParam("language"))
	pi.Sorts = []string{"vote_average.desc", "vote_average.asc", "first_air_date.desc", "first_air_date.asc", "popularity.desc", "popularity.asc"}

	limit := tmdb.ResultsPerPage * tmdb.PagesAtOnce
	pi.Pages = all/limit + 1
	if pi.Pages > 1000 {
		pi.Pages = 1000
	}
	for _, s := range shows {
		ii := new(ItemInfo)
		ii.ID = s.ID
		ii.Name = s.Name
		ii.OriginalName = s.OriginalName
		ii.Overview = s.Overview
		for _, g := range s.Genres {
			ii.Genres = append(ii.Genres, &tmdb.IDName{Name: g.Name, ID: g.ID})
		}
		if s.FirstAirDate != "" {
			ii.Year = s.FirstAirDate[:4]
		}
		ii.Poster = s.PosterPath
		ii.Backdrop = s.BackdropPath
		for _, i := range s.Images.Posters {
			ii.AllArts = append(ii.AllArts, i.FilePath)
		}
		for _, i := range s.Images.Backdrops {
			ii.AllArts = append(ii.AllArts, i.FilePath)
		}
		for _, i := range s.Images.Stills {
			ii.AllArts = append(ii.AllArts, i.FilePath)
		}
		ii.Seasons = s.NumberOfSeasons
		ii.Episodes = s.NumberOfEpisodes
		pi.Items = append(pi.Items, ii)
	}
	return pi
}
