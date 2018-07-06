package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"server/fdb"
	"server/fdb/tmdb"

	"github.com/labstack/echo"
)

func initSearch(e *echo.Echo) {
	e.GET("/search", searchPage)
	e.POST("/search/request", searchRequest)
	e.POST("/search/genres", genresRequest)
}

func searchPage(c echo.Context) error {
	return c.Render(http.StatusOK, "searchPage", nil)
}

type SearchRequest struct {
	Name         string
	Type         int
	Page         int
	HideWTorrent bool
	Filter       *tmdb.Filter `json:",omitempty"`
	SearchMovie  bool
	SearchTV     bool
}

func searchRequest(c echo.Context) error {
	jreq, err := getJsReqSearch(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	search := 0
	if jreq.SearchMovie && !jreq.SearchTV {
		search = 1
	} else if !jreq.SearchMovie && jreq.SearchTV {
		search = 2
	}

	var sResp *fdb.SearchResponce
	switch jreq.Type {
	case 1:
		sResp, err = fdb.NowWatching(jreq.Page, search)
	case 2:
		sResp, err = fdb.SearchByFilter(jreq.Page, jreq.Filter, search)
	default:
		if jreq.Name != "" {
			sResp, err = fdb.SearchByName(jreq.Page, jreq.Name, search)
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "Empty name")
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if jreq.HideWTorrent {
		var filtered = make([]*fdb.Movie, 0)
		for _, m := range sResp.Movies {
			if len(m.Torrents) > 0 {
				filtered = append(filtered, m)
			}
		}
		sResp.Movies = filtered
		return c.JSON(http.StatusOK, sResp)
	}
	return c.JSON(http.StatusOK, sResp)
}

func genresRequest(c echo.Context) error {
	ret := fdb.GetGenres()
	return c.JSON(http.StatusOK, ret)
}

func getJsReqSearch(c echo.Context) (*SearchRequest, error) {
	buf, _ := ioutil.ReadAll(c.Request().Body)
	jsstr := string(buf)
	decoder := json.NewDecoder(bytes.NewBufferString(jsstr))
	js := new(SearchRequest)
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
