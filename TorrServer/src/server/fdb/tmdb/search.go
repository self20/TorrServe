package tmdb

import (
	"fmt"

	"github.com/ryanbradynd05/go-tmdb"
)

func searchMovie(page int, name string) (*tmdb.MovieSearchResults, error) {
	var opt = make(map[string]string)
	opt["language"] = "ru"
	opt["page"] = fmt.Sprint(page)

	res, err := tm.SearchMovie(name, opt)
	if err != nil {
		return nil, err
	}
	fixLinks(res.Results)
	return res, err
}

func nowPlayingMovie(page int) (*tmdb.MovieDatedResults, error) {
	var opt = make(map[string]string)
	opt["language"] = "ru"
	opt["page"] = fmt.Sprint(page)

	res, err := tm.GetMovieNowPlaying(opt)
	if err != nil {
		return nil, err
	}
	fixLinks(res.Results)
	return res, nil
}

func discoverMovie(page int, filter *Filter) (*tmdb.MoviePagedResults, error) {
	var opt = filter.GetOptions()
	opt["page"] = fmt.Sprint(page)
	res, err := tm.DiscoverMovie(opt)
	if err != nil {
		return nil, err
	}
	fixLinks(res.Results)
	return res, nil
}

func searchTv(page int, name string) (*tmdb.TvPagedResults, error) {
	var opt = make(map[string]string)
	opt["language"] = "ru"
	opt["page"] = fmt.Sprint(page)

	res, err := tm.SearchTv(name, opt)
	if err != nil {
		return nil, err
	}

	ret := new(tmdb.TvPagedResults)
	ret.Page = res.Page
	ret.TotalPages = res.TotalPages
	ret.TotalResults = res.TotalResults
	for _, r := range res.Results {
		tvs := tmdb.TvShort{
			BackdropPath:  r.BackdropPath,
			ID:            r.ID,
			OriginalName:  r.OriginalName,
			FirstAirDate:  r.FirstAirDate,
			OriginCountry: r.OriginCountry,
			PosterPath:    r.PosterPath,
			Popularity:    r.Popularity,
			Name:          r.Name,
			VoteAverage:   r.VoteAverage,
			VoteCount:     r.VoteCount,
			GenreIds:      r.GenreIds,
		}
		ret.Results = append(ret.Results, tvs)
	}
	fixLinksTv(ret.Results)
	return ret, err
}

func nowPlayingTv(page int) (*tmdb.TvPagedResults, error) {
	var opt = make(map[string]string)
	opt["language"] = "ru"
	opt["page"] = fmt.Sprint(page)

	res, err := tm.GetTvOnTheAir(opt)
	if err != nil {
		return nil, err
	}
	fixLinksTv(res.Results)
	return res, nil
}

func discoverTv(page int, filter *Filter) (*tmdb.TvPagedResults, error) {
	var opt = filter.GetTvOptions()
	opt["page"] = fmt.Sprint(page)

	res, err := tm.DiscoverTV(opt)
	if err != nil {
		return nil, err
	}
	fixLinksTv(res.Results)
	return res, nil
}

func fixLinks(list []tmdb.MovieShort) {
	wbs := tmCfg.Images.BackdropSizes[2]
	wps := tmCfg.Images.PosterSizes[2]
	for i := 0; i < len(list); i++ {
		if list[i].BackdropPath != "" {
			list[i].BackdropPath = tmCfg.Images.BaseURL + wbs + list[i].BackdropPath
		}
		if list[i].PosterPath != "" {
			list[i].PosterPath = tmCfg.Images.BaseURL + wps + list[i].PosterPath
		}
	}
}

func fixLinksTv(list []tmdb.TvShort) {
	wbs := tmCfg.Images.BackdropSizes[2]
	wps := tmCfg.Images.PosterSizes[len(tmCfg.Images.PosterSizes)-1]
	for i := 0; i < len(list); i++ {
		if list[i].BackdropPath != "" {
			list[i].BackdropPath = tmCfg.Images.BaseURL + wbs + list[i].BackdropPath
		}
		if list[i].PosterPath != "" {
			list[i].PosterPath = tmCfg.Images.BaseURL + wps + list[i].PosterPath
		}
	}
}

//func removeDuplicateMovies(elements []tmdb.MovieShort) []tmdb.MovieShort {
//	encountered := map[tmdb.MovieShort]bool{}
//	var result []tmdb.MovieShort
//
//	for v := range elements {
//		if !encountered[elements[v]] {
//			encountered[elements[v]] = true
//			result = append(result, elements[v])
//		}
//	}
//	return result
//}

//
//func removeDuplicateTv(elements []tmdb.TvShort) []tmdb.TvShort {
//	encountered := make(map[tmdb.TvShort]bool, 0)
//	var result []tmdb.TvShort
//
//	for v := range elements {
//		if !encountered[elements[v]] {
//			encountered[elements[v]] = true
//			result = append(result, elements[v])
//		}
//	}
//	return result
//}
