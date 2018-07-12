package tmdb

import (
	"fmt"

	"github.com/jmcvetta/napping"
)

// GetSeason ...
func GetSeason(showID int, seasonNumber int, language string) *Season {
	var season *Season
	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":            apiKey,
			"append_to_response": "credits,images,videos,external_ids",
			"language":           language,
		}.AsUrlValues()
		resp, err := napping.Get(
			fmt.Sprintf("%stv/%d/season/%d", tmdbEndpoint, showID, seasonNumber),
			&urlValues,
			&season,
			nil,
		)
		if err != nil {
			fmt.Println(err.Error())
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded getting season %d of show %d, cooling down...\n", seasonNumber, showID)
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Printf("Bad status getting season %d of show %d: %d\n", seasonNumber, showID, resp.Status())
			return ErrHTTP
		}

		return nil
	})

	if season == nil {
		return nil
	}

	season.EpisodeCount = len(season.Episodes)

	// Fix for shows that have translations but return empty strings
	// for episode names and overviews.
	// We detect if episodes have their name filled, and if not re-query
	// with no language set.
	// See https://github.com/scakemyer/plugin.video.quasar/issues/249
	if season.EpisodeCount > 0 {
		for index := 0; index < season.EpisodeCount && index < len(season.Episodes); index++ {
			if season.Episodes[index] != nil && season.Episodes[index].Name == "" {
				season.Episodes[index] = GetEpisode(showID, seasonNumber, index+1, "")
			}
		}
	}

	return season
}

func (seasons SeasonList) Len() int           { return len(seasons) }
func (seasons SeasonList) Swap(i, j int)      { seasons[i], seasons[j] = seasons[j], seasons[i] }
func (seasons SeasonList) Less(i, j int) bool { return seasons[i].Season < seasons[j].Season }
