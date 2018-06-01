package provider

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Options struct {
	BaseUrl string
	Mirrors []string
}

type Torrent struct {
	Name    string
	Magnet  string
	Size    string
	PeersUl int
	PeersDl int
}

type Provider interface {
	Search(names, origNames string) ([]*Torrent, error)
	FindMirror()
}

func readPage(url string) (string, int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:21.0) Gecko/20100101 Firefox/21.0")
	client.Timeout = time.Second * 5
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	if resp.StatusCode != 200 {
		return "", resp.StatusCode, fmt.Errorf("%s %d", resp.Status, resp.StatusCode)
	}
	defer resp.Body.Close()

	var body io.Reader

	body = resp.Body

	if strings.Contains(resp.Header.Get("Content-Type"), "1251") {
		body = transform.NewReader(resp.Body, charmap.Windows1251.NewDecoder())
	}

	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	return string(buf), resp.StatusCode, nil
}
