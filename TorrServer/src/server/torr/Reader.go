package torr

import (
	"fmt"
	"io"

	"server/settings"

	"github.com/anacrolix/torrent"
)

type Reader struct {
	torrent.Reader
	file          *torrent.File
	currReadahead int64
	offset        int64
	lastOff       int64
}

func NewReader(file *torrent.File) *Reader {
	tr := file.NewReader()
	r := new(Reader)
	r.Reader = tr
	r.currReadahead = 1
	r.file = file
	tr.SetReadahead(r.currReadahead)
	return r
}

func (r *Reader) Read(p []byte) (n int, err error) {
	select {
	case <-r.file.Torrent().Closed():
		fmt.Println("Error, read closed torrent")
		return 0, io.ErrUnexpectedEOF
	default:

	}
	if r.currReadahead == 1 {
		r.currReadahead = getReadahead()
		r.Reader.SetReadahead(r.currReadahead)
	}
	n, err = r.Reader.Read(p)
	r.offset += int64(n)
	return
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	select {
	case <-r.file.Torrent().Closed():
		fmt.Println("Error, seek closed torrent")
		return 0, io.ErrUnexpectedEOF
	default:

	}
	r.currReadahead = 1
	r.Reader.SetReadahead(r.currReadahead)
	var err error
	r.offset, err = r.Reader.Seek(offset, whence)
	fmt.Println("Seek:", r.offset, r.currReadahead)
	return r.offset, err
}

func getReadahead() int64 {
	rhp := settings.Get().ReadAhead
	if rhp < 1 || rhp > 100 {
		rhp = 33
	}
	readahead := settings.Get().CacheSize / 100 * int64(rhp)
	return readahead
}

func (r *Reader) SetReadahead(buff int64) {
	r.Reader.SetReadahead(buff)
}
