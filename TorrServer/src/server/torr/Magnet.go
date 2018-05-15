package torr

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"server/settings"
	"server/utils"

	"github.com/anacrolix/torrent/metainfo"
)

func GetMagnet(link string) (string, error) {
	url, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	mag := ""
	switch strings.ToLower(url.Scheme) {
	case "magnet":
		mag, err = checkMagnet(url)
	case "http", "https":
		mag, err = getMagFromHttp(url.String())
	default:
		mag, err = getMagFromFile(url.Path)
	}
	if err != nil {
		return "", err
	}

	switch settings.Get().RetrackersMode {
	case 1:
		mag = utils.AddRetracker(mag)
	case 2:
		mag = utils.RemoveRetracker(mag)
	}

	return mag, nil
}

func checkMagnet(link *url.URL) (string, error) {
	hashs := link.Query()["xt"]
	for _, hs := range hashs {
		if strings.Contains(strings.ToLower(hs), "urn:btih:") {
			hash := strings.TrimPrefix(strings.ToLower(hs), "urn:btih:")
			if len(hash) != 40 {
				return "", errors.New("Wrong magnet link, size of hash not 40: " + link.String())
			}
			match, err := regexp.MatchString("^[0-9a-fA-F]+$", hash)
			if err != nil {
				return "", err
			}
			if !match {
				return "", errors.New("Wrong magnet link")
			}
		}
	}
	return link.String(), nil
}

func getMagFromHttp(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := new(http.Client)
	client.Timeout = time.Duration(time.Second * 30)
	req.Header.Set("User-Agent", "DWL/1.1.1 (Torrent)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	minfo, err := metainfo.Load(resp.Body)
	if err != nil {
		return "", err
	}
	info, err := minfo.UnmarshalInfo()
	if err != nil {
		return "", err
	}
	return minfo.Magnet(info.Name, minfo.HashInfoBytes()).String(), nil
}

func getMagFromFile(path string) (string, error) {
	minfo, err := metainfo.LoadFromFile(path)
	if err != nil {
		return "", err
	}
	info, err := minfo.UnmarshalInfo()
	if err != nil {
		return "", err
	}
	return minfo.Magnet(info.Name, minfo.HashInfoBytes()).String(), nil
}
