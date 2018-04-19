package torrent

import (
	"fmt"
	"io"
	"time"

	"torrentserver/settings"

	"github.com/anacrolix/torrent"
)

type Reader struct {
	io.Reader
	io.Seeker
	io.Closer

	tor    *torrent.Torrent
	file   *torrent.File
	reader torrent.Reader

	hash   string
	path   string
	offset int64

	lastRead time.Time
	closed   bool
	readNow  bool

	piecesLength int64
	pieceCurrent int
}

func NewReader(t *torrent.Torrent, f *torrent.File) *Reader {
	r := new(Reader)

	reader := f.NewReader()
	reader.SetReadahead(int64(settings.Get().PreloadBufferSize))

	r.tor = t
	r.file = f
	r.path = f.Path()
	r.hash = t.InfoHash().HexString()
	r.reader = reader
	r.piecesLength = r.tor.Info().PieceLength
	r.lastRead = time.Now().Add(time.Minute)
	r.readNow = true
	return r
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	if r.tor == nil {
		return 0, io.ErrUnexpectedEOF
	}

	r.readNow = true
	off, err := r.reader.Seek(offset, whence)
	r.readNow = false
	r.offset = off
	fmt.Println("Seek", r.offset, ", piece:", r.GetCurrentPiece())
	r.tor.PieceStateRuns()
	r.lastRead = time.Now().Add(time.Minute)
	return off, err
}

func (r *Reader) Read(p []byte) (n int, err error) {
	//if closed torrent return
	select {
	case <-r.tor.Closed():
		return 0, io.ErrUnexpectedEOF
	default:
	}

	r.readNow = true
	n, err = r.reader.Read(p)
	r.offset += int64(n)
	r.lastRead = time.Now()
	r.readNow = false

	readedPiece := r.GetCurrentPiece()
	if readedPiece != r.pieceCurrent {
		r.pieceCurrent = readedPiece
		storage.GetCache(r.hash).CurrentRead(readedPiece)
		fmt.Println("Read", r.offset, ", piece:", r.pieceCurrent)
	}
	return n, err
}

func (r *Reader) Close() error {
	r.reader.Close()
	r.closed = true
	r.readNow = false
	return nil
}

func (r *Reader) IsClosed() bool {
	return r.closed
}

func (r *Reader) IsExpired() bool {
	return !r.readNow && r.lastRead.Add(time.Minute).Before(time.Now())
}

func (r *Reader) GetCurrentPiece() int {
	return int((r.file.Offset() + r.offset) / r.piecesLength)
}
