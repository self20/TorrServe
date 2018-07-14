package parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type YoHoHo struct {
	url string
}

func NewYHH() *YoHoHo {
	p := new(YoHoHo)
	p.url = "https://4h0y.yohoho.cc"
	return p
}

func (p *YoHoHo) Search(findString string) ([]*Torrent, error) {
	fmt.Println("Yohoho finding:", findString)
	return p.findTorrents(findString)
}

func (p *YoHoHo) FindMirror() {}

func (p *YoHoHo) findTorrents(name string) ([]*Torrent, error) {
	t := &url.URL{Path: name}
	ename := t.String()
	url := fmt.Sprintf("%s/?title=%s", p.url, ename)
	body, _, err := readPage(url)
	if err != nil {
		return nil, err
	}

	return p.parse(body)
}

func (p *YoHoHo) parse(buf string) ([]*Torrent, error) {
	buf = strings.Replace(buf, "\n", " ", -1)
	reg, err := regexp.Compile(`<span class="td-btn" onclick="window\.location\.href =.+?'(magnet:\?.+?)';">(.+?)<\/span>.+?<div.+?>(.+?)<`)
	if err != nil {
		return nil, err
	}
	src := reg.FindAllStringSubmatch(buf, -1)
	if len(src) > 0 {
		tors := make([]*Torrent, 0)
		for _, t := range src {
			t[3] = strings.Replace(t[3], "&nbsp;", " ", -1)
			tor := new(Torrent)
			tor.Magnet = t[1]
			tor.Name = t[2]
			tor.Size = t[3]
			tor.PeersDl = -1
			tor.PeersUl = -1
			tors = append(tors, tor)
		}
		return tors, nil
	}
	return nil, nil
}
