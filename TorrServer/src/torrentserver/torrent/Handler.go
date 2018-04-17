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
			deletedHandler := 0
			for hi := range h.Handlers {
				handle := h.Handlers[hi-deletedHandler]
				deleted := 0
				for i := range handle.Readers {
					j := i - deleted
					if handle.Readers[j].lastRead.Add(time.Minute).Before(time.Now()) {
						handle.Readers[j].Close()
						handle.Readers = append(handle.Readers[:j], handle.Readers[j+1:]...)
						deleted++
						fmt.Println("Remove reader:", handle.Torrent.Name(), len(handle.Readers))
					}
				}
				if len(handle.Readers) == 0 {
					fmt.Println("Remove torrent:", handle.Torrent.Name())
					handle.Torrent.Drop()
					j := hi - deletedHandler
					h.Handlers = append(h.Handlers[:j], h.Handlers[j+1:]...)
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
