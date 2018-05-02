package torrent

import (
	"fmt"
	"io"
	"sync"

	"torrentserver/settings"

	"github.com/anacrolix/torrent"
)

var (
	count = 0
	mu    sync.Mutex
)

type Reader struct {
	io.Reader
	io.Seeker
	io.Closer

	index int

	tor    *torrent.Torrent
	file   *torrent.File
	reader torrent.Reader

	hash   string
	path   string
	offset int64

	closed bool

	piecesLength int64
	pieceCurrent int
}

func NewReader(t *torrent.Torrent, f *torrent.File) *Reader {
	r := new(Reader)

	reader := f.NewReader()

	mu.Lock()
	count++
	r.index = count
	mu.Unlock()

	r.tor = t
	r.file = f
	r.path = f.Path()
	r.hash = t.InfoHash().HexString()
	r.reader = reader
	r.piecesLength = r.tor.Info().PieceLength
	if settings.Get().IsElementumCache {
		reader.SetReadahead(int64(float64(settings.Get().CacheSize) * 0.33))
	} else {
		reader.SetReadahead(int64(float64(settings.Get().CacheSize) * 0.33))
	}
	return r
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	if r.tor == nil {
		return 0, io.ErrUnexpectedEOF
	}

	off, err := r.reader.Seek(offset, whence)
	r.offset = off
	fmt.Println("Seek", r.index, r.offset, ", piece:", r.GetCurrentPiece(), "/", r.GetCountPieces())
	return off, err
}

func (r *Reader) Read(p []byte) (n int, err error) {
	//if closed torrent return
	select {
	case <-r.tor.Closed():
		return 0, io.ErrUnexpectedEOF
	default:
	}

	n, err = r.reader.Read(p)
	r.offset += int64(n)

	readedPiece := r.GetCurrentPiece()
	if readedPiece != r.pieceCurrent {
		r.pieceCurrent = readedPiece
		if settings.Get().IsElementumCache {
			storage.GetElCache(r.hash).CurrentRead(readedPiece)
		} else {
			storage.GetCache(r.hash).CurrentRead(readedPiece)
		}
		fmt.Println("Read", r.index, r.offset, ", piece:", r.pieceCurrent, "/", r.GetCountPieces())
	}
	return n, err
}

func (r *Reader) Close() error {
	r.reader.Close()
	r.closed = true
	return nil
}

func (r *Reader) IsClosed() bool {
	return r.closed
}

func (r *Reader) GetCurrentPiece() int {
	return int((r.file.Offset() + r.offset) / r.piecesLength)
}

func (r *Reader) GetCountPieces() int {
	return int((r.file.Offset() + r.file.Length()) / r.piecesLength)
}
