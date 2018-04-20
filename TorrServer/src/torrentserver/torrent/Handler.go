package torrent

import (
	"fmt"
	"time"

	"torrentserver/db"

	"github.com/anacrolix/sync"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

var ()

type Handler struct {
	Handlers []*Handle

	mu       sync.Mutex
	watching bool
}

type Handle struct {
	expired time.Time
	Torrent *torrent.Torrent
	Readers []*Reader
}

func NewHandler() *Handler {
	h := new(Handler)
	h.Handlers = make([]*Handle, 0)
	return h
}

func (h *Handler) watch() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.watching {
		return
	}
	h.watching = true
	go func() {
		for h.watching && len(h.Handlers) > 0 {
			h.mu.Lock()

			for handleIndex := 0; handleIndex < len(h.Handlers); handleIndex++ {
				handle := h.Handlers[handleIndex]
				for readerIndex := 0; readerIndex < len(handle.Readers); readerIndex++ {
					if handle.Readers[readerIndex].IsClosed() {
						if h.removeReader(handle, readerIndex) {
							if readerIndex > 0 {
								readerIndex--
							}
							fmt.Println("Remove reader:", handle.Torrent.Name(), len(handle.Readers))
						}
					}
				}
				if len(handle.Readers) == 0 && time.Now().After(handle.expired) {
					if h.removeTorrent(handleIndex) {
						if handleIndex > 0 {
							handleIndex--
						}
						fmt.Println("Remove torrent:", handle.Torrent.Name())
					}
				}
			}
			h.mu.Unlock()
			time.Sleep(time.Second)
		}
		h.watching = false
	}()
}

func (h *Handler) NewReader(torrent *db.Torrent, filename string) (*Reader, error) {
	torr, err := getTorrent(torrent)
	if err != nil {
		return nil, err
	}
	for _, f := range torr.Files() {
		if f.Path() == filename {
			reader := NewReader(torr, f)
			h.addReader(torr, reader)
			h.watch()
			return reader, nil
		}
	}
	return nil, fmt.Errorf("File in torrent not found: %v/%v", torr.InfoHash().HexString(), filename)
}

func (h *Handler) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, h := range h.Handlers {
		h.Torrent.Drop()
	}
}

func (h *Handler) addReader(tor *torrent.Torrent, reader *Reader) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, h := range h.Handlers {
		if h.Torrent == tor {
			h.Readers = append(h.Readers, reader)
			fmt.Println("Add reader:", tor.Name(), len(h.Readers))
			return
		}
	}
	handl := new(Handle)
	handl.Torrent = tor
	handl.Readers = append(handl.Readers, reader)
	h.Handlers = append(h.Handlers, handl)
	fmt.Println("Add reader:", tor.Name())
}

func (h *Handler) removeTorrent(torrentIndex int) bool {
	if torrentIndex >= 0 && torrentIndex < len(h.Handlers) {
		h.Handlers[torrentIndex].Torrent.Drop()
		h.Handlers = append(h.Handlers[:torrentIndex], h.Handlers[torrentIndex+1:]...)
		return true
	}
	return false
}

func (h *Handler) removeReader(handle *Handle, readerIndex int) bool {
	if readerIndex >= 0 && readerIndex < len(handle.Readers) {
		handle.Readers[readerIndex].Close()
		handle.Readers = append(handle.Readers[:readerIndex], handle.Readers[readerIndex+1:]...)
		if len(handle.Readers) == 0 {
			handle.expired = time.Now().Add(time.Minute)
		}
		return true
	}
	return false
}

func getTorrent(tordb *db.Torrent) (*torrent.Torrent, error) {
	hash := metainfo.NewHashFromHex(tordb.Hash)
	if tor, ok := client.Torrent(hash); ok {
		return tor, nil
	}
	tor, err := client.AddMagnet(tordb.Magnet)
	if err != nil {
		return nil, err
	}
	err = GotInfo(tor)
	if err != nil {
		return nil, err
	}
	return tor, nil
}
