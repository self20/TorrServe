package fdb

import (
	"fmt"
	"sync"

	"server/fdb/provider"
	"server/fdb/tmdb"
	"server/utils"
)

type Movie struct {
	Id          int
	Title       string
	OrigTitle   string
	Date        string
	BackdropUrl string
	PosterUrl   string
	Overview    string
	IsTv        bool
	Torrents    []*provider.Torrent
}

type SearchResponce struct {
	Movies []*Movie
	Pages  int
}

var (
	providers []provider.Provider
)

func Init() {
	//providers = append(providers, provider.NewRutor(provider.Options{}))
	providers = append(providers, provider.NewYHH(provider.Options{}))
	//providers = append(providers, provider.NewTParser(provider.Options{}))

	var wa sync.WaitGroup
	wa.Add(len(providers))
	for _, p := range providers {
		go func() {
			p.FindMirror()
			wa.Done()
		}()
	}
	wa.Wait()
}

func GetGenres() []struct {
	ID   int
	Name string
} {
	return tmdb.GetGenres()
}

func SearchByName(page int, name string) (*SearchResponce, error) {
	var err error
	resp := new(SearchResponce)
	list := make([]*Movie, 0)
	fmt.Println("Search movies")
	{
		res, er := tmdb.SearchMovie(page, name)
		if err != nil {
			err = er
		} else {
			resp.Pages = res.TotalPages
			for _, m := range res.Results {
				if m.Title == m.OriginalTitle && !utils.IsCyrillic(m.OriginalTitle) {
					continue
				}
				mov := new(Movie)
				mov.Id = m.ID
				mov.Title = m.Title
				mov.OrigTitle = m.OriginalTitle
				mov.Date = m.ReleaseDate
				mov.BackdropUrl = m.BackdropPath
				mov.PosterUrl = m.PosterPath
				mov.Overview = m.Overview
				list = append(list, mov)
			}
		}
	}
	{
		fmt.Println("Search tv")
		res, er := tmdb.SearchTv(page, name)
		if err != nil {
			err = er
		} else {
			if res.TotalPages > resp.Pages {
				resp.Pages = res.TotalPages
			}
			for _, m := range res.Results {
				if m.Name == m.OriginalName {
					continue
				}
				mov := new(Movie)
				mov.Id = m.ID
				mov.Title = m.Name
				mov.OrigTitle = m.OriginalName
				mov.Date = m.FirstAirDate
				mov.BackdropUrl = m.BackdropPath
				mov.PosterUrl = m.PosterPath
				mov.Overview = m.Overview
				mov.IsTv = true
				list = append(list, mov)
			}
		}
	}

	if len(list) > 0 {
		findTorrents(list)
	}
	resp.Movies = list
	return resp, err
}

func NowWatching(page int) (*SearchResponce, error) {
	var err error
	resp := new(SearchResponce)
	list := make([]*Movie, 0)
	fmt.Println("Search now watching movies")
	{
		res, er := tmdb.NowPlayingMovie(page)
		if err != nil {
			err = er
		} else {
			resp.Pages = res.TotalPages
			for _, m := range res.Results {
				if m.Title == m.OriginalTitle {
					continue
				}
				mov := new(Movie)
				mov.Id = m.ID
				mov.Title = m.Title
				mov.OrigTitle = m.OriginalTitle
				mov.Date = m.ReleaseDate
				mov.BackdropUrl = m.BackdropPath
				mov.PosterUrl = m.PosterPath
				mov.Overview = m.Overview
				list = append(list, mov)
			}
		}
	}
	{
		fmt.Println("Search now watching tv")
		res, er := tmdb.NowPlayingTv(page)
		if err != nil {
			err = er
		} else {
			if res.TotalPages > resp.Pages {
				resp.Pages = res.TotalPages
			}
			for _, m := range res.Results {
				if m.Name == m.OriginalName {
					continue
				}
				mov := new(Movie)
				mov.Id = m.ID
				mov.Title = m.Name
				mov.OrigTitle = m.OriginalName
				mov.Date = m.FirstAirDate
				mov.BackdropUrl = m.BackdropPath
				mov.PosterUrl = m.PosterPath
				mov.Overview = m.Overview
				mov.IsTv = true
				list = append(list, mov)
			}
		}
	}

	if len(list) > 0 {
		findTorrents(list)
	}
	resp.Movies = list
	return resp, err
}

func SearchByFilter(page int, filter *tmdb.Filter) (*SearchResponce, error) {
	var err error
	resp := new(SearchResponce)
	list := make([]*Movie, 0)
	fmt.Println("Search filter movies")
	{
		res, er := tmdb.FilterMovie(page, filter)
		if err != nil {
			err = er
		} else {
			resp.Pages = res.TotalPages
			for _, m := range res.Results {
				if m.Title == m.OriginalTitle {
					continue
				}
				mov := new(Movie)
				mov.Id = m.ID
				mov.Title = m.Title
				mov.OrigTitle = m.OriginalTitle
				mov.Date = m.ReleaseDate
				mov.BackdropUrl = m.BackdropPath
				mov.PosterUrl = m.PosterPath
				mov.Overview = m.Overview
				list = append(list, mov)
			}
		}
	}
	{
		fmt.Println("Search filter tv")
		res, er := tmdb.FilterTv(page, filter)
		if err != nil {
			err = er
		} else {
			if res.TotalPages > resp.Pages {
				resp.Pages = res.TotalPages
			}
			for _, m := range res.Results {
				if m.Name == m.OriginalName {
					continue
				}
				mov := new(Movie)
				mov.Id = m.ID
				mov.Title = m.Name
				mov.OrigTitle = m.OriginalName
				mov.Date = m.FirstAirDate
				mov.BackdropUrl = m.BackdropPath
				mov.PosterUrl = m.PosterPath
				mov.Overview = m.Overview
				mov.IsTv = true
				list = append(list, mov)
			}
		}
	}
	if len(list) > 0 {
		findTorrents(list)
	}
	resp.Movies = list
	return resp, err
}

func findTorrents(movies []*Movie) {
	utils.ParallelFor(0, len(movies), func(i int) {
		movie := movies[i]
		var torList []*provider.Torrent
		for _, p := range providers {
			res, err := p.Search(movie.Title, movie.OrigTitle)
			if err == nil {
				torList = append(torList, res...)
			}
		}
		movie.Torrents = torList
	})
}
