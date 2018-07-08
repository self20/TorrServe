package provider

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type YoHoHo struct {
	opt     Options
	antiBan bool
}

func NewYHH(opt Options) *YoHoHo {
	p := new(YoHoHo)
	p.opt = opt
	if p.opt.BaseUrl == "" {
		p.opt.BaseUrl = "https://4h0y.yohoho.cc"
	}
	return p
}

func (p *YoHoHo) Search(findString string) ([]*Torrent, error) {
	fmt.Println("Find torrents:", findString)
	return p.findTorrents(findString)
}

func (p *YoHoHo) FindMirror() {

}

func (p *YoHoHo) findTorrents(name string) ([]*Torrent, error) {
	t := &url.URL{Path: name}
	ename := t.String()
	url := fmt.Sprintf("%s/?title=%s", p.opt.BaseUrl, ename)
	body, _, err := readPage(url)
	if err != nil {
		return nil, err
	}
	if p.antiBan {
		sp := rand.Intn(int(time.Millisecond * 600))
		time.Sleep(time.Millisecond*500 + time.Duration(sp))
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
