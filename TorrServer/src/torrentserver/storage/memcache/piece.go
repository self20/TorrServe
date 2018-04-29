package memcache

import (
	"errors"
	"io"
	"time"

	"github.com/anacrolix/torrent/storage"
)

type Piece struct {
	storage.PieceImpl

	Id     int
	Hash   string
	Length int64

	complete bool
	readed   bool
	accessed time.Time
	buffer   []byte
}

func (p *Piece) WriteAt(b []byte, off int64) (n int, err error) {
	if p.buffer == nil {
		p.buffer = make([]byte, p.Length)
	}
	n = copy(p.buffer[off:], b[:])
	p.accessed = time.Now().Add(time.Second * 5)
	return
}

func (p *Piece) ReadAt(b []byte, off int64) (n int, err error) {
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
	p.buffer = nil
	p.complete = false
}

func (p *Piece) Stat() ItemState {
	itm := ItemState{
		Id:         p.Id,
		Hash:       p.Hash,
		Accessed:   p.accessed,
		Completed:  p.complete,
		BufferSize: len(p.buffer),
	}
	return itm
}
