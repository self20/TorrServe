package torrent

import (
	"fmt"
	"sync"
	"time"

	"torrentserver/settings"

	"github.com/anacrolix/torrent"
)

type Preloader struct {
	offset int64
	length int64
	file   string

	preload bool
	mu      sync.Mutex
}

type PreloadStat struct {
	PreloadFile   string
	IsPreload     bool
	PreloadOffset int64
	PreloadLength int64
}

func NewPreloader() *Preloader {
	p := new(Preloader)
	return p
}

func (p *Preloader) Preload(file *torrent.File) {
	if settings.Get().PreloadBufferSize == 0 {
		return
	}

	go func() {
		p.mu.Lock()
		if p.preload {
			p.mu.Unlock()
			return
		}
		p.preload = true
		defer func() { p.preload = false }()
		p.mu.Unlock()

		pieceLength := file.Torrent().Info().PieceLength

		p.offset = int64(0)
		p.length = int64(settings.Get().PreloadBufferSize)
		p.file = file.Path()

		ps := int((file.Offset() + p.offset) / pieceLength)
		pe := int((file.Offset() + p.length) / pieceLength)

		fmt.Println("Starting preload peieces:", ps, "-", pe)

		reader := file.NewReader()
		defer reader.Close()
		reader.SetReadahead(p.length - p.offset)

		buf := make([]byte, 65536)
		update := int64(0)
		for p.offset < p.length && p.preload {
			if update > pieceLength {
				fmt.Println("Preloaded:", p.offset, "/", p.length)
				update = 0
				reader.SetReadahead(p.length - p.offset)
			}
			readed, err := reader.Read(buf)
			if err != nil {
				fmt.Println("Error read preload:", err)
				return
			}
			p.offset += int64(readed)
			update += int64(readed)
		}
	}()
	time.Sleep(time.Millisecond * 500)
}

func (p *Preloader) Stat() PreloadStat {
	return PreloadStat{
		p.file,
		p.preload,
		p.offset,
		p.length,
	}
}

func (p *Preloader) Stop() {
	p.preload = false
}
