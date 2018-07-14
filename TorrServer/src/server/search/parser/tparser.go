package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"server/utils"
)

type TParser struct {
}

var jsNum = map[string]string{
	"2": "1",
	//"4":  "2",
	"6": "3",
	//"8":  "4",
	"9":  "5",
	"10": "5",
}

func NewTParser() *TParser {
	p := new(TParser)
	return p
}

func (p *TParser) FindMirror() {}

func (p *TParser) Search(findString string) ([]*Torrent, error) {
	fmt.Println("TParser finding:", findString)
	urls := make([]string, 0)

	for k, v := range jsNum {
		nurl := fmt.Sprintf("http://js%v.tparser.org/js%v/%v.tor.php?callback=one&jsonpx=%v&s=1",
			v, v, k, url.PathEscape(findString))
		urls = append(urls, nurl)
	}

	return p.findTorrents(urls)
}

func (p *TParser) findTorrents(urls []string) ([]*Torrent, error) {
	list := make([]*Torrent, 0)
	var err error
	var mu sync.Mutex
	for _, u := range urls {
		body, _, er := readPage(u)
		if er != nil {
			err = er
		} else {
			//fmt.Println("Parse:", u)
			tors, er := p.parse(body)
			if err != nil {
				err = er
			} else {
				mu.Lock()
				list = append(list, tors...)
				mu.Unlock()
			}
		}
	}

	if len(list) == 0 && err != nil {
		return nil, err
	}

	return list, nil
}

type tparserJS struct {
	Sr []struct {
		Name string `json:"name"`
		Size string `json:"size"`
		S    string `json:"s"`
		L    string `json:"l"`
		Link string `json:"link"`
		T    string `json:"t"`
		Sk   string `json:"sk"`
		Img  string `json:"img"`
		K    string `json:"k"`
		Z    string `json:"z"`
		D    string `json:"d"`
	} `json:"sr"`
}

func (p *TParser) parse(buf string) ([]*Torrent, error) {
	buf = buf[4 : len(buf)-1]
	buf = strings.Replace(buf, "'", "\"", -1)
	js := new(tparserJS)
	err := json.Unmarshal([]byte(buf), js)
	if err != nil {
		return nil, err
	}

	torrList := make([]*Torrent, 0)
	var er error
	utils.ParallelFor(0, len(js.Sr), func(i int) {
		sr := js.Sr[i]
		if sr.Z == "1" {
			torr := new(Torrent)
			torr.Name = strings.Replace(strings.Replace(sr.Name, "<b>", "", -1), "</b>", "", -1)
			torr.Size = sr.Size + " " + sr.T
			torr.PeersDl, _ = strconv.Atoi(sr.L)
			torr.PeersUl, _ = strconv.Atoi(sr.S)
			mag, err := p.getMagnet(sr.Img, sr.D)
			if err != nil {
				er = err
			} else {
				torr.Magnet = mag
				torrList = append(torrList, torr)
			}
		}
	})
	if len(torrList) == 0 && er != nil {
		return nil, er
	}
	return torrList, nil
}

func (p *TParser) getMagnet(img, d string) (string, error) {
	link := "http://tparser.org/magnet.php?t=" + strconv.Itoa(len(img)) + img + d
	//fmt.Println("Get magnet:", link)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:21.0) Gecko/20100101 Firefox/21.0")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 300 && resp.StatusCode < 309 {
		mag := resp.Header.Get("Location")
		if mag == "" {
			mag = resp.Header.Get("Content-Location")
		}
		if mag != "" {
			//fmt.Println("Found magnet:", mag)
			return mag, nil
		}
	}
	return "", errors.New("Magnet not found")
}
