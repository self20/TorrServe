package torrent

import (
	"fmt"
	"sync"

	"torrentserver/settings"

	"github.com/anacrolix/torrent"
)

type Preloader struct {
	file *torrent.File

	offset int64
	length int64

	preload bool
	mu      sync.Mutex
}

type PreloadStat struct {
	PreloadFile   string
	IsPreload     bool
	PreloadOffset int64
	PreloadLength int64
}

func NewPreloader(file *torrent.File) *Preloader {
	if file == nil {
		return nil
	}
	p := new(Preloader)
	p.file = file
	return p
}

func (p *Preloader) Preload() {
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

		pieceLength := p.file.Torrent().Info().PieceLength

		p.offset = int64(0)
		p.length = int64(settings.Get().PreloadBufferSize)

		ps := int((p.file.Offset() + p.offset) / pieceLength)
		pe := int((p.file.Offset() + p.length) / pieceLength)

		fmt.Println("Starting preload peieces:", ps, "-", pe)

		reader := p.file.NewReader()
		reader.SetReadahead(int64(float64(settings.Get().PreloadBufferSize) * 0.33))
		defer reader.Close()

		buf := make([]byte, 65536)
		for p.offset < p.length && p.preload {
			readed, err := reader.Read(buf)
			if err != nil {
				fmt.Println("Error read preload:", err)
				return
			}
			p.offset += int64(readed)
		}
	}()
}

func (p *Preloader) Stat() PreloadStat {
	return PreloadStat{
		p.file.Path(),
		p.preload,
		p.offset,
		p.length,
	}
}

func (p *Preloader) Stop() {
	p.preload = false
}
