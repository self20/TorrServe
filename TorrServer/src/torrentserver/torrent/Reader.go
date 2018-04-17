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

	hash     string
	path     string
	offset   int64
	lastRead time.Time

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
	return r
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	if r.tor == nil {
		return 0, io.ErrUnexpectedEOF
	}

	off, err := r.reader.Seek(offset, whence)
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

	n, err = r.reader.Read(p)
	r.offset += int64(n)
	r.lastRead = time.Now()

	readedPiece := r.GetCurrentPiece()
	if readedPiece != r.pieceCurrent {
		fmt.Println("Read", r.offset, ", piece:", r.pieceCurrent)
		r.pieceCurrent = readedPiece
		storage.GetCache(r.hash).CurrentRead(readedPiece)
	}
	return n, err
}

func (r *Reader) Close() error {
	r.reader.Close()
	return nil
}

func (r *Reader) GetCurrentPiece() int {
	return int((r.file.Offset() + r.offset) / r.piecesLength)
}
