package tmdb

import (
	"fmt"
)

type Filter struct {
	SortAsc       bool
	SortBy        string
	DateLte       string
	DateGte       string
	WithGenres    []int
	WithoutGenres []int
}

func (f *Filter) GetOptions() map[string]string {
	var opt = make(map[string]string)
	opt["language"] = "ru"

	by := ".asc"
	if !f.SortAsc {
		by = ".desc"
	}
	switch f.SortBy {
	case "popularity", "release_date", "revenue", "primary_release_date", "original_title", "vote_average", "vote_count":
		opt["sort_by"] = f.SortBy + by
	}
	if f.DateLte != "" {
		//2014-10-22
		opt["release_date.lte"] = f.DateLte
	}
	if f.DateGte != "" {
		opt["release_date.gte"] = f.DateGte
	}

	if len(f.WithGenres) >= 1 {
		opt["with_genres"] = fmt.Sprint(f.WithGenres[0])
	}
	for i := 1; i < len(f.WithGenres); i++ {
		opt["with_genres"] += "," + fmt.Sprint(f.WithGenres[i])
	}
	if len(f.WithoutGenres) >= 1 {
		opt["without_genres"] = fmt.Sprint(f.WithoutGenres[0])
	}
	for i := 1; i < len(f.WithoutGenres); i++ {
		opt["without_genres"] += "," + fmt.Sprint(f.WithoutGenres[i])
	}
	return opt
}

func (f *Filter) GetTvOptions() map[string]string {
	var opt = make(map[string]string)
	opt["language"] = "ru"

	by := ".asc"
	if !f.SortAsc {
		by = ".desc"
	}
	switch f.SortBy {

	case "vote_average", "first_air_date", "popularity":
		opt["sort_by"] = f.SortBy + by
	}
	if f.DateLte != "" {
		//2014-10-22
		opt["first_air_date.lte"] = f.DateLte
	}
	if f.DateGte != "" {
		opt["first_air_date.gte"] = f.DateGte
	}

	if len(f.WithGenres) >= 1 {
		opt["with_genres"] = fmt.Sprint(f.WithGenres[0])
	}
	for i := 1; i < len(f.WithGenres); i++ {
		opt["with_genres"] += "," + fmt.Sprint(f.WithGenres[i])
	}
	return opt
}
