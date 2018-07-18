package tmdb

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmcvetta/napping"
)

// GetShowImages ...
func GetShowImages(showID int) *Images {
	var images *Images
	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":                apiKey,
			"include_image_language": fmt.Sprintf("%s,en,null", "ru"),
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"tv/"+strconv.Itoa(showID)+"/images",
			&urlValues,
			&images,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded getting images for %d, cooling down...\n", showID)
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Printf("Bad status getting images for %d: %d\n", showID, resp.Status())
			return ErrHTTP
		}

		return nil
	})
	return images
}

// GetShowByID ...
func GetShowByID(tmdbID string, language string) *Show {
	id, _ := strconv.Atoi(tmdbID)
	return GetShow(id, language)
}

// GetShow ...
func GetShow(showID int, language string) (show *Show) {
	if showID == 0 {
		return
	}
	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":            apiKey,
			"append_to_response": "credits,images,alternative_titles,translations,external_ids",
			"language":           language,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"tv/"+strconv.Itoa(showID),
			&urlValues,
			&show,
			nil,
		)
		if err != nil {
			switch e := err.(type) {
			case *json.UnmarshalTypeError:
				fmt.Printf("UnmarshalTypeError: Value[%s] Type[%v] Offset[%d] for %d\n", e.Value, e.Type, e.Offset, showID)
			case *json.InvalidUnmarshalError:
				fmt.Printf("InvalidUnmarshalError: Type[%v]\n", e.Type)
			default:
				fmt.Println(err)
			}
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded getting show %d, cooling down...\n", showID)
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			message := fmt.Sprintf("Bad status getting show for %d: %d\n", showID, resp.Status())
			fmt.Println(message)
			return ErrHTTP
		}

		return nil
	})
	if show == nil {
		return nil
	}

	switch t := show.RawPopularity.(type) {
	case string:
		if popularity, err := strconv.ParseFloat(t, 64); err == nil {
			show.Popularity = popularity
		}
	case float64:
		show.Popularity = t
	}

	return show
}

// GetShows ...
func GetShows(showIds []int, language string) Shows {
	var wg sync.WaitGroup
	shows := make(Shows, len(showIds))
	wg.Add(len(showIds))
	for i, showID := range showIds {
		go func(i int, showId int) {
			defer wg.Done()
			shows[i] = GetShow(showId, language)
		}(i, showID)
	}
	wg.Wait()
	return shows
}

// SearchShows ...
func SearchShows(query string, language string, page int) (Shows, int) {
	var results EntityList
	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key": apiKey,
			"query":   query,
			"page":    strconv.Itoa(page),
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"search/tv",
			&urlValues,
			&results,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded searching shows for %s, cooling down...\n", query)
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Printf("Bad status searching shows: %d\n", resp.Status())
			return ErrHTTP
		}

		return nil
	})
	tmdbIds := make([]int, 0, len(results.Results))
	for _, entity := range results.Results {
		tmdbIds = append(tmdbIds, entity.ID)
	}
	return GetShows(tmdbIds, language), results.TotalResults
}

func listShows(endpoint string, params napping.Params, page int) (Shows, int) {
	params["api_key"] = apiKey
	totalResults := -1

	limit := ResultsPerPage * PagesAtOnce
	//pageGroup := (page-1)*ResultsPerPage/limit + 1

	shows := make(Shows, limit)

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
					fmt.Printf("Rate limit exceeded while listing shows from %s, cooling down...\n", endpoint)
					rl.CoolDown(resp.HttpResponse().Header)
					return ErrExceeded
				} else if resp.Status() != 200 {
					message := fmt.Sprintf("Bad status while listing shows: %d\n", resp.Status())
					fmt.Println(message)
					return ErrHTTP
				}

				return nil
			})
			if results != nil {
				totalResults = results.TotalResults
				var wgItems sync.WaitGroup
				wgItems.Add(len(results.Results))
				for s, show := range results.Results {
					if show == nil {
						wgItems.Done()
						continue
					}

					go func(i int, tmdbId int) {
						defer wgItems.Done()
						shows[i] = GetShow(tmdbId, params["language"])
					}(p*ResultsPerPage+s, show.ID)
				}
				wgItems.Wait()
			}
		}(p)
	}
	wg.Wait()

	return shows, totalResults
}

// PopularShows ...
func PopularShows(params DiscoverFilters, language string, page int) (Shows, int) {
	var p napping.Params
	if params.Genre != "" {
		p = napping.Params{
			"language":           language,
			"sort_by":            "popularity.desc",
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
			"with_genres":        params.Genre,
		}
	} else if params.Country != "" {
		p = napping.Params{
			"language":           language,
			"sort_by":            "popularity.desc",
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
			"region":             params.Country,
		}
	} else if params.Language != "" {
		p = napping.Params{
			"language":               language,
			"sort_by":                "popularity.desc",
			"first_air_date.lte":     time.Now().UTC().Format("2006-01-02"),
			"with_original_language": params.Language,
		}
	} else {
		p = napping.Params{
			"language":           language,
			"sort_by":            "popularity.desc",
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
		}
	}

	return listShows("discover/tv", p, page)
}

func DiscoverShows(params map[string]string, page int) (Shows, int) {
	if _, ok := params["first_air_date.lte"]; !ok {
		params["first_air_date.lte"] = time.Now().UTC().Format("2006-01-02")
	}

	if _, ok := params["language"]; !ok {
		params["language"] = "ru"
	}

	return listShows("discover/tv", params, page)
}

// RecentShows ...
func RecentShows(params DiscoverFilters, language string, page int) (Shows, int) {
	var p napping.Params
	if params.Genre != "" {
		p = napping.Params{
			"language":           language,
			"sort_by":            "first_air_date.desc",
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
			"with_genres":        params.Genre,
		}
	} else if params.Country != "" {
		p = napping.Params{
			"language":           language,
			"sort_by":            "first_air_date.desc",
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
			"region":             params.Country,
		}
	} else if params.Language != "" {
		p = napping.Params{
			"language":               language,
			"sort_by":                "first_air_date.desc",
			"first_air_date.lte":     time.Now().UTC().Format("2006-01-02"),
			"with_original_language": params.Language,
		}
	} else {
		p = napping.Params{
			"language":           language,
			"sort_by":            "first_air_date.desc",
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
		}
	}

	return listShows("discover/tv", p, page)
}

// RecentEpisodes ...
func RecentEpisodes(params DiscoverFilters, language string, page int) (Shows, int) {
	var p napping.Params

	if params.Genre != "" {
		p = napping.Params{
			"language":           language,
			"air_date.gte":       time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02"),
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
			"with_genres":        params.Genre,
		}
	} else if params.Country != "" {
		p = napping.Params{
			"language":           language,
			"air_date.gte":       time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02"),
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
			"region":             params.Country,
		}
	} else if params.Language != "" {
		p = napping.Params{
			"language":               language,
			"air_date.gte":           time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02"),
			"first_air_date.lte":     time.Now().UTC().Format("2006-01-02"),
			"with_original_language": params.Language,
		}
	} else {
		p = napping.Params{
			"language":           language,
			"air_date.gte":       time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02"),
			"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
		}
	}

	return listShows("discover/tv", p, page)
}

// TopRatedShows ...
func TopRatedShows(genre string, language string, page int) (Shows, int) {
	return listShows("tv/top_rated", napping.Params{"language": language}, page)
}

// MostVotedShows ...
func MostVotedShows(genre string, language string, page int) (Shows, int) {
	return listShows("discover/tv", napping.Params{
		"language":           language,
		"sort_by":            "vote_count.desc",
		"first_air_date.lte": time.Now().UTC().Format("2006-01-02"),
		"with_genres":        genre,
	}, page)
}

// GetTVGenres ...
func GetTVGenres(language string) []*Genre {
	genres := GenreList{}

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":  apiKey,
			"language": language,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"genre/tv/list",
			&urlValues,
			&genres,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded getting TV genres, cooling down...\n")
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Printf("Bad status getting TV genres: %d\n", resp.Status())
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
