package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Rutor struct {
	url string
}

var rtmirrors = []string{
	"http://top-tor.org",
	"http://free-ru.org",
	"http://zerkalo-rutor.org",
	"http://free-rutor.org",
	"http://fast-bit.org",

	//Не официальные зеркала
	//"http://srutor.org",
	//"http://nerutor.org",
}

func NewRutor() *Rutor {
	p := new(Rutor)
	p.url = "http://rutor.info"
	p.FindMirror()
	return p
}

func (p *Rutor) Search(findString string) ([]*Torrent, error) {
	fmt.Println("Rutor finding:", findString)
	return p.findTorrents(findString)
}

func (p *Rutor) FindMirror() {
	_, code, err := readPage(p.url)
	if code == 200 && err == nil {
		return
	}
	fmt.Println("Find mirror rutor:")
	for _, m := range rtmirrors {
		fmt.Println("Check:", m)
		_, code, err := readPage(m)
		if code == 200 && err == nil {
			fmt.Println("Find:", m)
			p.url = m
			return
		}
	}
}

func (p *Rutor) findTorrents(name string) ([]*Torrent, error) {
	url := fmt.Sprintf("%s/search/0/0/100/2/%s", p.url, name)
	body, _, err := readPage(url)
	if err != nil {
		return nil, err
	}
	return p.parse(body)
}

func (p *Rutor) parse(buf string) ([]*Torrent, error) {
	buf = strings.Replace(buf, "\n", " ", -1)
	reg, err := regexp.Compile(`tr class="(?:gai|tum)".+?"(magnet.+?)".+?<a href=".+?">(.+?)<\/a>.+?<td align="right">?(\d+?\.\d+?.+?)<\/td.+?alt="S" \/>(.+?)<\/span>.+?<span class="red">(.+?)<\/span>`)
	if err != nil {
		return nil, err
	}
	src := reg.FindAllStringSubmatch(buf, -1)
	if len(src) > 0 {
		tors := make([]*Torrent, 0)
		for _, t := range src {
			t[3] = strings.TrimSpace(strings.Replace(t[3], "&nbsp;", " ", -1))
			t[4] = strings.TrimSpace(strings.Replace(t[4], "&nbsp;", " ", -1))
			t[5] = strings.TrimSpace(strings.Replace(t[5], "&nbsp;", " ", -1))

			tor := new(Torrent)
			tor.Magnet = t[1]
			tor.Name = t[2]
			tor.Size = t[3]
			tor.PeersUl, err = strconv.Atoi(t[4])
			if err != nil {
				tor.PeersUl = -1
			}
			tor.PeersDl, err = strconv.Atoi(t[5])
			if err != nil {
				tor.PeersDl = -1
			}
			tors = append(tors, tor)
		}
		return tors, nil
	}
	return nil, nil
}
