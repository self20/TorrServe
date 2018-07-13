package tmdb

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmcvetta/napping"
)

// ByPopularity ...
type ByPopularity Movies

func (a ByPopularity) Len() int           { return len(a) }
func (a ByPopularity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPopularity) Less(i, j int) bool { return a[i].Popularity < a[j].Popularity }

// GetImages ...
func GetImages(movieID int) *Images {
	var images *Images

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":                apiKey,
			"include_image_language": "ru,en,null",
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"movie/"+strconv.Itoa(movieID)+"/images",
			&urlValues,
			&images,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded getting images for %d, cooling down...\n", movieID)
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Printf("Bad status getting images for %d: %d\n", movieID, resp.Status())
			return ErrHTTP
		}

		return nil
	})
	return images
}

// GetMovie ...
func GetMovie(tmdbID int, language string) *Movie {
	return GetMovieByID(strconv.Itoa(tmdbID), language)
}

// GetMovieByID ...
func GetMovieByID(movieID string, language string) *Movie {
	var movie *Movie

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":            apiKey,
			"append_to_response": "credits,images,alternative_titles,translations,external_ids,trailers,release_dates",
			"language":           language,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"movie/"+movieID,
			&urlValues,
			&movie,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded getting movie %s, cooling down...\n", movieID)
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Printf("Bad status getting movie %s: %d\n", movieID, resp.Status())
			return ErrHTTP
		}

		return nil
	})

	if movie == nil {
		return nil
	}
	switch t := movie.RawPopularity.(type) {
	case string:
		popularity, _ := strconv.ParseFloat(t, 64)
		movie.Popularity = popularity
	case float64:
		movie.Popularity = t
	}
	return movie
}

// GetMovies ...
func GetMovies(tmdbIds []int, language string) Movies {
	var wg sync.WaitGroup
	movies := make(Movies, len(tmdbIds))
	wg.Add(len(tmdbIds))
	for i, tmdbID := range tmdbIds {
		go func(i int, tmdbId int) {
			defer wg.Done()
			movies[i] = GetMovie(tmdbId, language)
		}(i, tmdbID)
	}
	wg.Wait()
	return movies
}

// GetMovieGenres ...
func GetMovieGenres(language string) []*Genre {
	genres := GenreList{}

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":  apiKey,
			"language": language,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"genre/movie/list",
			&urlValues,
			&genres,
			nil,
		)

		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Println("Rate limit exceeded getting genres, cooling down...")
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			message := fmt.Sprintf("Bad status getting movie genres: %d", resp.Status())
			fmt.Println(message)
			return ErrHTTP
		}

		return nil
	})
	if genres.Genres != nil && len(genres.Genres) > 0 {
		for _, i := range genres.Genres {
			i.Name = strings.Title(i.Name)
		}

		sort.Slice(genres.Genres, func(i, j int) bool {
			return genres.Genres[i].Name < genres.Genres[j].Name
		})
	}
	return genres.Genres
}

// SearchMovies ...
func SearchMovies(query string, language string, page int) (Movies, int) {
	var results EntityList

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key": apiKey,
			"query":   query,
			"page":    strconv.Itoa(page),
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"search/movie",
			&urlValues,
			&results,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded searching movies with %s\n", query)
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Printf("Bad status searching movies: %d\n", resp.Status())
			return ErrHTTP
		}

		return nil
	})
	tmdbIds := make([]int, 0, len(results.Results))
	for _, movie := range results.Results {
		tmdbIds = append(tmdbIds, movie.ID)
	}
	return GetMovies(tmdbIds, language), results.TotalResults
}

// GetIMDBList ...
func GetIMDBList(listID string, language string, page int) (movies Movies, totalResults int) {
	var results *List
	totalResults = -1
	limit := ResultsPerPage * PagesAtOnce

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key": apiKey,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"list/"+listID,
			&urlValues,
			&results,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Println("Rate limit exceeded getting IMDb list, cooling down...")
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			message := fmt.Sprintf("Bad status getting IMDb list: %d", resp.Status())
			fmt.Println(message + fmt.Sprintf(" (%s)", listID))
			return ErrHTTP
		}

		return nil
	})
	tmdbIds := make([]int, 0)
	for i, movie := range results.Items {
		if i >= limit {
			break
		}
		tmdbIds = append(tmdbIds, movie.ID)
	}
	movies = GetMovies(tmdbIds, language)
	return
}

func listMovies(endpoint string, params napping.Params, page int) (Movies, int) {
	params["api_key"] = apiKey
	totalResults := -1

	limit := ResultsPerPage * PagesAtOnce
	//pageGroup := (page-1)*ResultsPerPage/limit + 1

	movies := make(Movies, limit)

	wg := sync.WaitGroup{}
	for p := 0; p < PagesAtOnce; p++ {
		wg.Add(1)
		//currentPage := (pageGroup-1)*ResultsPerPage + p + 1
		go func(p int) {
			defer wg.Done()
			var results *EntityList
			pageParams := napping.Params{
				"page": strconv.Itoa(page),
			}
			for k, v := range params {
				pageParams[k] = v
			}
			urlParams := pageParams.AsUrlValues()
			rl.Call(func() error {
				resp, err := napping.Get(
					tmdbEndpoint+endpoint,
					&urlParams,
					&results,
					nil,
				)
				if err != nil {
					fmt.Println(err)
				} else if resp.Status() == 429 {
					fmt.Printf("Rate limit exceeded listing movies from %s, cooling down...\n", endpoint)
					rl.CoolDown(resp.HttpResponse().Header)
					return ErrExceeded
				} else if resp.Status() != 200 {
					fmt.Printf("Bad status while listing movies from %s: %d\n", endpoint, resp.Status())
					return ErrHTTP
				}

				return nil
			})
			if results != nil {
				totalResults = results.TotalResults
				var wgItems sync.WaitGroup
				wgItems.Add(len(results.Results))
				for m, movie := range results.Results {
					if movie == nil {
						wgItems.Done()
						continue
					}

					go func(i int, tmdbId int) {
						defer wgItems.Done()
						movies[i] = GetMovie(tmdbId, params["language"])
					}(p*ResultsPerPage+m, movie.ID)
				}
				wgItems.Wait()
			}
		}(p)
	}
	wg.Wait()
	return movies, totalResults
}

// PopularMovies ...
func PopularMovies(params DiscoverFilters, language string, page int) (Movies, int) {
	var p napping.Params
	if params.Genre != "" {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "popularity.desc",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
			"with_genres":              params.Genre,
		}
	} else if params.Country != "" {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "popularity.desc",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
			"region":                   params.Country,
		}
	} else if params.Language != "" {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "popularity.desc",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
			"with_original_language":   params.Language,
		}
	} else {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "popularity.desc",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
		}
	}

	return listMovies("discover/movie", p, page)
}

func DiscoverMovies(params map[string]string, page int) (Movies, int) {
	//if _, ok := params["vote_count.gte"]; !ok {
	//	params["vote_count.gte"] = "10"
	//}

	//if _, ok := params["primary_release_date.lte"]; !ok {
	//	params["primary_release_date.lte"] = time.Now().UTC().Format("2006-01-02")
	//}

	if _, ok := params["language"]; !ok {
		params["language"] = "ru"
	}

	return listMovies("discover/movie", params, page)
}

// RecentMovies ...
func RecentMovies(params DiscoverFilters, language string, page int) (Movies, int) {
	var p napping.Params
	if params.Genre != "" {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "primary_release_date.desc",
			"vote_count.gte":           "10",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
			"with_genres":              params.Genre,
		}
	} else if params.Country != "" {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "primary_release_date.desc",
			"vote_count.gte":           "10",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
			"region":                   params.Country,
		}
	} else if params.Language != "" {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "primary_release_date.desc",
			"vote_count.gte":           "10",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
			"with_original_language":   params.Language,
		}
	} else {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "primary_release_date.desc",
			"vote_count.gte":           "10",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
		}
	}

	return listMovies("discover/movie", p, page)
}

// TopRatedMovies ...
func TopRatedMovies(genre string, language string, page int) (Movies, int) {
	var p napping.Params
	if genre == "" {
		p = napping.Params{
			"language": language,
		}
	} else {
		p = napping.Params{
			"language":    language,
			"with_genres": genre,
		}
	}
	return listMovies("movie/top_rated", p, page)
}

// MostVotedMovies ...
func MostVotedMovies(genre string, language string, page int) (Movies, int) {
	var p napping.Params
	if genre == "" {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "vote_count.desc",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
		}
	} else {
		p = napping.Params{
			"language":                 language,
			"sort_by":                  "vote_count.desc",
			"primary_release_date.lte": time.Now().UTC().Format("2006-01-02"),
			"with_genres":              genre,
		}
	}
	return listMovies("discover/movie", p, page)
}
