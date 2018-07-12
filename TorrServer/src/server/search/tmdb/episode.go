package tmdb

import (
	"fmt"

	"github.com/jmcvetta/napping"
)

// GetEpisode ...
func GetEpisode(showID int, seasonNumber int, episodeNumber int, language string) *Episode {
	var episode *Episode

	rl.Call(func() error {
		urlValues := napping.Params{
			"api_key":            apiKey,
			"append_to_response": "credits,images,videos,external_ids",
			"language":           language,
		}.AsUrlValues()
		resp, err := napping.Get(
			fmt.Sprintf("%stv/%d/season/%d/episode/%d", tmdbEndpoint, showID, seasonNumber, episodeNumber),
			&urlValues,
			&episode,
			nil,
		)
		if err != nil {
			fmt.Println(err.Error())
		} else if resp.Status() == 429 {
			fmt.Printf("Rate limit exceeded getting S%02dE%02d of %d, cooling down...\n", seasonNumber, episodeNumber, showID)
			rl.CoolDown(resp.HttpResponse().Header)
			return ErrExceeded
		} else if resp.Status() != 200 {
			fmt.Printf("Bad status getting S%02dE%02d of %d: %d\n", seasonNumber, episodeNumber, showID, resp.Status())
			return ErrHTTP
		}

		return nil
	})

	return episode
}
