package memcache

import (
	"errors"
	"io"
	"sync"
	"time"

	"torrentserver/storage/state"

	"github.com/anacrolix/torrent/storage"
)

type Piece struct {
	storage.PieceImpl

	Id     int
	Hash   string
	Length int64
	Size   int64

	complete bool
	readed   bool
	accessed time.Time
	buffer   []byte

	mu    sync.RWMutex
	cache *Cache
}

func (p *Piece) WriteAt(b []byte, off int64) (n int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.buffer == nil {
		p.buffer = p.cache.bufferPull.GetBuffer()
	}
	n = copy(p.buffer[off:], b[:])
	p.Size += int64(n)
	p.accessed = time.Now()

	p.cache.cleanPieces()
	return
}

func (p *Piece) ReadAt(b []byte, off int64) (n int, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	size := len(b)
	if size+int(off) > len(p.buffer) {
		size = len(p.buffer) - int(off)
		if size < 0 {
			size = 0
		}
	}
	if len(p.buffer) < int(off) || len(p.buffer) < int(off)+size {
		return 0, io.ErrUnexpectedEOF
	}
	n = copy(b, p.buffer[int(off) : int(off)+size][:])
	p.accessed = time.Now()
	if int(off)+size >= len(p.buffer) {
		p.readed = true
	}

	p.cache.cleanPieces()
	return n, nil
}

func (p *Piece) MarkComplete() error {
	if len(p.buffer) == 0 {
		return errors.New("piece is not complete")
	}
	p.complete = true
	return nil
}

func (p *Piece) MarkNotComplete() error {
	p.complete = false
	return nil
}

func (p *Piece) Completion() storage.Completion {
	return storage.Completion{
		Complete: p.complete && len(p.buffer) > 0,
		Ok:       true,
	}
}

func (p *Piece) Release() {
	p.cache.bufferPull.PutBuffer(p.buffer)
	p.buffer = nil
	p.Size = 0
	p.complete = false
}

func (p *Piece) Stat() state.ItemState {
	itm := state.ItemState{
		Id:         p.Id,
		Hash:       p.Hash,
		Accessed:   p.accessed,
		Completed:  p.complete,
		BufferSize: p.Size,
	}
	return itm
}
