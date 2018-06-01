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
	Name   string
	Type   int
	Page   int
	Filter *tmdb.Filter `json:",omitempty"`
}

func searchRequest(c echo.Context) error {
	jreq, err := getJsReqSearch(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var movies []*fdb.Movie
	switch jreq.Type {
	case 1:
		movies, err = fdb.NowWatching(jreq.Page)
	case 2:
		movies, err = fdb.SearchByFilter(jreq.Page, jreq.Filter)
	default:
		if jreq.Name != "" {
			movies, err = fdb.SearchByName(jreq.Page, jreq.Name)
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "Empty name")
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, movies)
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
