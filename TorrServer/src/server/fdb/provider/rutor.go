package provider

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Rutor struct {
	opt     Options
	antiBan bool
	mu      sync.Mutex
}

var mirrors = []string{
	"http://top-tor.org",
	"http://free-ru.org",
	"http://zerkalo-rutor.org",
	"http://free-rutor.org",
	"http://fast-bit.org",

	//Не официальные зеркала
	"http://srutor.org",
	"http://nerutor.org",
}

func NewRutor(opt Options) *Rutor {
	p := new(Rutor)
	p.opt = opt
	if p.opt.BaseUrl == "" {
		p.opt.BaseUrl = "http://rutor.info"
	}
	if len(p.opt.Mirrors) == 0 {
		for _, m := range mirrors {
			p.opt.Mirrors = append(p.opt.Mirrors, m)
		}
	}
	return p
}

func (p *Rutor) Search(name, oname string) ([]*Torrent, error) {
	if p.antiBan {
		p.mu.Lock()
		defer p.mu.Unlock()
	}
	fmt.Println("Find torrents:", name, "/", oname)
	tors, err := p.findTorrents(name)
	if len(tors) == 0 {
		tors, err = p.findTorrents(oname)
	}

	if p.antiBan {
		sp := rand.Intn(int(time.Millisecond * 600))
		time.Sleep(time.Millisecond*500 + time.Duration(sp))
	}
	return tors, err
}

func (p *Rutor) FindMirror() {
	_, code, err := readPage(p.opt.BaseUrl)
	if code == 200 && err == nil {
		p.antiBan = true
		return
	}
	fmt.Println("Find mirror rutor:")
	for i, m := range p.opt.Mirrors {
		fmt.Println("Check:", m)
		_, code, err := readPage(m)
		if code == 200 && err == nil {
			fmt.Println("Find:", m)
			p.opt.BaseUrl = m
			p.antiBan = i < 5
			return
		}
	}
}

func (p *Rutor) findTorrents(name string) ([]*Torrent, error) {
	url := fmt.Sprintf("%s/search/%s", p.opt.BaseUrl, name)
	body, _, err := readPage(url)
	if err != nil {
		return nil, err
	}
	return p.parse(body)
}

func (p *Rutor) parse(buf string) ([]*Torrent, error) {
	buf = strings.Replace(buf, "\n", " ", -1)
	reg, err := regexp.Compile(`"(magnet:\?.+?)".+?<a.+?>(.+?)<\/a>.+<td align="right">(\d+?\.?\d+?.+?).+?<span class="green">(?:<img.+?>)?(.+?)<\/span>.+?<span class="red">(.+?)<\/span>`)
	if err != nil {
		return nil, err
	}
	src := reg.FindAllStringSubmatch(buf, -1)
	if len(src) > 0 {
		tors := make([]*Torrent, 0)
		for _, t := range src {

			t[4] = strings.Replace(t[4], "&nbsp;", " ", -1)
			t[5] = strings.Replace(t[5], "&nbsp;", " ", -1)

			tor := new(Torrent)
			tor.Magnet = t[1]
			tor.Name = t[2]
			tor.Size = t[3]
			tor.PeersUl, err = strconv.Atoi(t[4])
			tor.PeersDl, err = strconv.Atoi(t[5])
			tors = append(tors, tor)
		}
		return tors, nil
	}
	return nil, nil
}
