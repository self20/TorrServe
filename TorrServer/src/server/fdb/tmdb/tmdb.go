package tmdb

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/ryanbradynd05/go-tmdb"
)

var (
	apiKeys = []string{
		"8cf43ad9c085135b9479ad5cf6bbcbda",
		"ae4bd1b6fce2a5648671bfc171d15ba4",
		"29a551a65eef108dd01b46e27eb0554a",
	}
	tm     *tmdb.TMDb
	tmCfg  *tmdb.Configuration
	genres *tmdb.Genre
)

func initDb() error {
	key := apiKeys[rand.Intn(len(apiKeys))]
	fmt.Println("Key:", key)
	tm = tmdb.Init(key)
	if tmCfg == nil {
		cfg, err := tm.GetConfiguration()
		tmCfg = cfg
		if err != nil {
			return err
		}
	}

	if genres == nil {
		var opt = make(map[string]string)
		opt["language"] = "ru"
		gen, err := tm.GetMovieGenres(opt)
		genres = gen
		if err != nil {
			return err
		}
		gen, err = tm.GetTvGenres(opt)
		genres.Genres = append(genres.Genres, gen.Genres...)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetGenres() []struct {
	ID   uint32
	Name string
} {
	initDb()
	return genres.Genres
}

func SearchMovie(page int, name string) (*tmdb.MovieSearchResults, error) {
	err := initDb()
	if err != nil {
		return nil, err
	}

	res, err := searchMovie(page, name)
	if err != nil {
		return nil, err
	}

	//res.Results = removeDuplicateMovies(res.Results)
	sort.Slice(res.Results, func(i, j int) bool {
		return res.Results[i].Popularity > res.Results[j].Popularity
	})

	return res, nil
}

func SearchTv(page int, name string) (*tmdb.TvPagedResults, error) {
	err := initDb()
	if err != nil {
		return nil, err
	}

	res, err := searchTv(page, name)
	if err != nil {
		return nil, err
	}

	sort.Slice(res.Results, func(i, j int) bool {
		return res.Results[i].Popularity > res.Results[j].Popularity
	})

	return res, nil
}

func NowPlayingMovie(page int) (*tmdb.MovieDatedResults, error) {
	err := initDb()
	if err != nil {
		return nil, err
	}

	res, err := nowPlayingMovie(page)
	if err != nil {
		return nil, err
	}
	sort.Slice(res.Results, func(i, j int) bool {
		return res.Results[i].Popularity > res.Results[j].Popularity
	})

	return res, nil
}

func NowPlayingTv(page int) (*tmdb.TvPagedResults, error) {
	err := initDb()
	if err != nil {
		return nil, err
	}

	res, err := nowPlayingTv(page)
	if err != nil {
		return nil, err
	}

	//res.Results = removeDuplicateTv(res.Results)
	sort.Slice(res.Results, func(i, j int) bool {
		return res.Results[i].Popularity > res.Results[j].Popularity
	})

	return res, nil
}

func FilterMovie(page int, filter *Filter) (*tmdb.MoviePagedResults, error) {
	err := initDb()
	if err != nil {
		return nil, err
	}

	res, err := discoverMovie(page, filter)

	return res, nil
}

func FilterTv(page int, filter *Filter) (*tmdb.TvPagedResults, error) {
	err := initDb()
	if err != nil {
		return nil, err
	}

	res, err := discoverTv(page, filter)

	return res, nil
}
