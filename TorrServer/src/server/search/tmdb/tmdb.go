package tmdb

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/jmcvetta/napping"
)

//go:generate msgp -o msgp.go -io=false -tests=false

const (
	// PagesAtOnce ...
	PagesAtOnce = 1
	// ResultsPerPage ...
	ResultsPerPage = 20
)

// Movies ...
type Movies []*Movie

// Shows ...
type Shows []*Show

// SeasonList ...
type SeasonList []*Season

// EpisodeList ...
type EpisodeList []*Episode

// Movie ...
type Movie struct {
	Entity

	IMDBId              string       `json:"imdb_id"`
	Overview            string       `json:"overview"`
	ProductionCompanies []*IDName    `json:"production_companies"`
	Runtime             int          `json:"runtime"`
	TagLine             string       `json:"tagline"`
	RawPopularity       interface{}  `json:"popularity"`
	Popularity          float64      `json:"-"`
	SpokenLanguages     []*Language  `json:"spoken_languages"`
	ExternalIDs         *ExternalIDs `json:"external_ids"`

	AlternativeTitles *struct {
		Titles []*AlternativeTitle `json:"titles"`
	} `json:"alternative_titles"`

	Translations *struct {
		Translations []*Language `json:"translations"`
	} `json:"translations"`

	Trailers *struct {
		Youtube []*Trailer `json:"youtube"`
	} `json:"trailers"`

	Credits *Credits `json:"credits,omitempty"`
	Images  *Images  `json:"images,omitempty"`

	ReleaseDates *ReleaseDatesResults `json:"release_dates"`
}

// Show ...
type Show struct {
	Entity

	EpisodeRunTime      []int        `json:"episode_run_time"`
	Genres              []*Genre     `json:"genres"`
	Homepage            string       `json:"homepage"`
	InProduction        bool         `json:"in_production"`
	FirstAirDate        string       `json:"first_air_date"`
	LastAirDate         string       `json:"last_air_date"`
	Networks            []*IDName    `json:"networks"`
	NumberOfEpisodes    int          `json:"number_of_episodes"`
	NumberOfSeasons     int          `json:"number_of_seasons"`
	OriginalName        string       `json:"original_name"`
	OriginCountry       []string     `json:"origin_country"`
	Overview            string       `json:"overview"`
	RawPopularity       interface{}  `json:"popularity"`
	Popularity          float64      `json:"-"`
	ProductionCompanies []*IDName    `json:"production_companies"`
	Status              string       `json:"status"`
	ExternalIDs         *ExternalIDs `json:"external_ids"`
	Translations        *struct {
		Translations []*Language `json:"translations"`
	} `json:"translations"`
	AlternativeTitles *struct {
		Titles []*AlternativeTitle `json:"results"`
	} `json:"alternative_titles"`

	Credits *Credits `json:"credits,omitempty"`
	Images  *Images  `json:"images,omitempty"`

	Seasons SeasonList `json:"seasons"`
}

// Season ...
type Season struct {
	ID           int          `json:"id"`
	Name         string       `json:"name,omitempty"`
	Season       int          `json:"season_number"`
	EpisodeCount int          `json:"episode_count,omitempty"`
	AirDate      string       `json:"air_date"`
	Poster       string       `json:"poster_path"`
	ExternalIDs  *ExternalIDs `json:"external_ids"`

	Episodes EpisodeList `json:"episodes"`
}

// Episode ...
type Episode struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Overview      string       `json:"overview"`
	AirDate       string       `json:"air_date"`
	SeasonNumber  int          `json:"season_number"`
	EpisodeNumber int          `json:"episode_number"`
	VoteAverage   float32      `json:"vote_average"`
	StillPath     string       `json:"still_path"`
	ExternalIDs   *ExternalIDs `json:"external_ids"`
}

// Entity ...
type Entity struct {
	IsAdult          bool      `json:"adult"`
	BackdropPath     string    `json:"backdrop_path"`
	ID               int       `json:"id"`
	Genres           []*IDName `json:"genres"`
	OriginalTitle    string    `json:"original_title,omitempty"`
	OriginalLanguage string    `json:"original_language,omitempty"`
	ReleaseDate      string    `json:"release_date"`
	FirstAirDate     string    `json:"first_air_date"`
	PosterPath       string    `json:"poster_path"`
	Title            string    `json:"title,omitempty"`
	VoteAverage      float32   `json:"vote_average"`
	VoteCount        int       `json:"vote_count"`
	OriginalName     string    `json:"original_name,omitempty"`
	Name             string    `json:"name,omitempty"`
}

// EntityList ...
type EntityList struct {
	Page         int       `json:"page"`
	Results      []*Entity `json:"results"`
	TotalPages   int       `json:"total_pages"`
	TotalResults int       `json:"total_results"`
}

// IDName ...
type IDName struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Genre ...
type Genre IDName

// GenreList ...
type GenreList struct {
	Genres []*Genre `json:"genres"`
}

// Country ...
type Country struct {
	Iso31661    string `json:"iso_3166_1"`
	EnglishName string `json:"english_name"`
}

// CountryList ...
type CountryList []*Country

// LanguageList ...
type LanguageList struct {
	Languages []*Language `json:"languages"`
}

// Image ...
type Image struct {
	FilePath string `json:"file_path"`
	Height   int    `json:"height"`
	Iso639_1 string `json:"iso_639_1"`
	Width    int    `json:"width"`
}

// Images ...
type Images struct {
	Backdrops []*Image `json:"backdrops"`
	Posters   []*Image `json:"posters"`
	Stills    []*Image `json:"stills"`
}

// Cast ...
type Cast struct {
	IDName
	CastID      int    `json:"cast_id"`
	Character   string `json:"character"`
	CreditID    string `json:"credit_id"`
	Order       int    `json:"order"`
	ProfilePath string `json:"profile_path"`
}

// Crew ...
type Crew struct {
	IDName
	CreditID    string `json:"credit_id"`
	Department  string `json:"department"`
	Job         string `json:"job"`
	ProfilePath string `json:"profile_path"`
}

// Credits ...
type Credits struct {
	Cast []*Cast `json:"cast"`
	Crew []*Crew `json:"crew"`
}

// ExternalIDs ...
type ExternalIDs struct {
	IMDBId      string      `json:"imdb_id"`
	FreeBaseID  string      `json:"freebase_id"`
	FreeBaseMID string      `json:"freebase_mid"`
	TVDBID      interface{} `json:"tvdb_id"`
}

// AlternativeTitle ...
type AlternativeTitle struct {
	Iso3166_1 string `json:"iso_3166_1"`
	Title     string `json:"title"`
}

// Language ...
type Language struct {
	Iso639_1    string `json:"iso_639_1"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name,omitempty"`
}

// Configuration ...
type Configuration struct {
	Images struct {
		BaseURL       string   `json:"base_url"`
		SecureBaseURL string   `json:"secure_base_url"`
		BackdropSizes []string `json:"backdrop_sizes"`
		LogoSizes     []string `json:"logo_sizes"`
		PosterSizes   []string `json:"poster_sizes"`
		ProfileSizes  []string `json:"profile_sizes"`
		StillSizes    []string `json:"still_sizes"`
	}
	ChangeKeys []string `json:"change_keys,omitempty"`
}

// FindResult ...
type FindResult struct {
	MovieResults     []*Entity `json:"movie_results"`
	PersonResults    []*Entity `json:"person_results"`
	TVResults        []*Entity `json:"tv_results"`
	TVEpisodeResults []*Entity `json:"tv_episode_results"`
	TVSeasonResults  []*Entity `json:"tv_season_results"`
}

// List ...
type List struct {
	CreatedBy     string    `json:"created_by"`
	Description   string    `json:"description"`
	FavoriteCount int       `json:"favorite_count"`
	ID            string    `json:"id"`
	ItemCount     int       `json:"item_count"`
	Iso639_1      string    `json:"iso_639_1"`
	Name          string    `json:"name"`
	PosterPath    string    `json:"poster_path"`
	Items         []*Entity `json:"items"`
}

// Trailer ...
type Trailer struct {
	Name   string `json:"name"`
	Size   string `json:"size"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

// ReleaseDatesResults ...
type ReleaseDatesResults struct {
	Results []*ReleaseDates `json:"results"`
}

// ReleaseDates ...
type ReleaseDates struct {
	Iso3166_1    string         `json:"iso_3166_1"`
	ReleaseDates []*ReleaseDate `json:"release_dates"`
}

// ReleaseDate ...
type ReleaseDate struct {
	Certification string `json:"certification"`
	Iso639_1      string `json:"iso_639_1"`
	Note          string `json:"note"`
	ReleaseDate   string `json:"release_date"`
	Type          int    `json:"type"`
}

// DiscoverFilters ...
type DiscoverFilters struct {
	Genre    string
	Country  string
	Language string
}

const (
	tmdbEndpoint            = "https://api.themoviedb.org/3/"
	imageEndpoint           = "http://image.tmdb.org/t/p/"
	burstRate               = 40
	burstTime               = 10 * time.Second
	simultaneousConnections = 20
)

var (
	apiKeys = []string{
		"8cf43ad9c085135b9479ad5cf6bbcbda",
		"ae4bd1b6fce2a5648671bfc171d15ba4",
		"29a551a65eef108dd01b46e27eb0554a",
	}
	apiKey = apiKeys[rand.Intn(len(apiKeys))]
	// WarmingUp ...
	WarmingUp = true
)

var rl = NewRateLimiter(burstRate, burstTime, simultaneousConnections)

// CheckAPIKey ...
func CheckAPIKey() {
	result := false
	for index := len(apiKeys) - 1; index >= 0; index-- {
		result = tmdbCheck(apiKey)
		if result {
			break
		} else {
			if apiKey == apiKeys[index] {
				apiKeys = append(apiKeys[:index], apiKeys[index+1:]...)
			}
			if len(apiKeys) > 0 {
				apiKey = apiKeys[rand.Intn(len(apiKeys))]
			} else {
				result = false
				break
			}
		}
	}
	if result == false {
		fmt.Println("No valid TMDB API key found")
	}
}

func tmdbCheck(key string) bool {
	var result *Entity

	urlValues := napping.Params{
		"api_key": key,
	}.AsUrlValues()

	resp, err := napping.Get(
		tmdbEndpoint+"movie/550",
		&urlValues,
		&result,
		nil,
	)

	if err != nil {
		return false
	} else if resp.Status() != 200 {
		return false
	}

	return true
}

// ImageURL ...
func ImageURL(uri string, size string) string {
	return imageEndpoint + size + uri
}

// ListEntities ...
// TODO Unused...
func ListEntities(endpoint string, params napping.Params) []*Entity {
	var wg sync.WaitGroup
	resultsPerPage := 28
	entities := make([]*Entity, PagesAtOnce*resultsPerPage)
	params["api_key"] = apiKey
	params["language"] = "ru"

	wg.Add(PagesAtOnce)
	for i := 0; i < PagesAtOnce; i++ {
		go func(page int) {
			defer wg.Done()
			var tmp *EntityList
			tmpParams := napping.Params{
				"page": strconv.Itoa(page),
			}
			for k, v := range params {
				tmpParams[k] = v
			}
			urlValues := tmpParams.AsUrlValues()
			rl.Call(func() error {
				resp, err := napping.Get(
					tmdbEndpoint+endpoint,
					&urlValues,
					&tmp,
					nil,
				)
				if err != nil {
					fmt.Println(err.Error())
				} else if resp.Status() != 200 {
					message := fmt.Sprintf("Bad status listing entities: %d", resp.Status())
					fmt.Println(message)
				}

				return nil
			})
			for i, entity := range tmp.Results {
				entities[page*resultsPerPage+i] = entity
			}
		}(i)
	}
	wg.Wait()

	return entities
}

// Find ...
func Find(externalID string, externalSource string) *FindResult {
	var result *FindResult

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":         apiKey,
			"external_source": externalSource,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"find/"+externalID,
			&urlValues,
			&result,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Println("Rate limit exceeded finding tmdb item, cooling down...")
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Println("Bad status finding tmdb item:", resp.Status())
		}

		return nil
	})

	return result
}

// GetCountries ...
func GetCountries(language string) []*Country {
	countries := CountryList{}

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key": apiKey,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"configuration/countries",
			&urlValues,
			&countries,
			nil,
		)

		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Println("Rate limit exceeded getting countries, cooling down...")
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			message := fmt.Sprintf("Bad status getting countries: %d", resp.Status())
			fmt.Println(message)
			return ErrHTTP
		}

		return nil
	})
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].EnglishName < countries[j].EnglishName
	})
	return countries
}

// GetLanguages ...
func GetLanguages(language string) []*Language {
	languages := []*Language{}

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key": apiKey,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"configuration/languages",
			&urlValues,
			&languages,
			nil,
		)

		if err != nil {
			fmt.Println(err)
		} else if resp.Status() == 429 {
			fmt.Println("Rate limit exceeded getting languages, cooling down...")
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			message := fmt.Sprintf("Bad status getting languages: %d", resp.Status())
			fmt.Println(message)
			return ErrHTTP
		}

		return nil
	})
	for _, l := range languages {
		if l.Name == "" {
			l.Name = l.EnglishName
		}
	}

	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Name < languages[j].Name
	})
	return languages
}

func GetConfig() Configuration {
	config := Configuration{}
	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key": apiKey,
		}.AsUrlValues()
		resp, err := napping.Get(
			tmdbEndpoint+"configuration",
			&urlValues,
			&config,
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

	return config
}
