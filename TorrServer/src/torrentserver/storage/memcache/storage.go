package memcache

import (
	"fmt"
	"sort"
	"sync"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

type Storage struct {
	storage.ClientImpl

	caches   map[string]*Cache
	capacity int
	mu       sync.Mutex
}

func NewStorage(capacity int) *Storage {
	stor := new(Storage)
	stor.capacity = capacity
	stor.caches = make(map[string]*Cache)
	return stor
}

func (s *Storage) OpenTorrent(info *metainfo.Info, infoHash metainfo.Hash) (storage.TorrentImpl, error) {
	fmt.Println("Open torrent", info.Name)
	ch := NewCache(s.capacity)
	ch.Init(info, infoHash)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.caches[infoHash.HexString()] = ch
	return ch, nil
}

func (s *Storage) CloseByHash(hash string) {
	if s.caches == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if ch, ok := s.caches[hash]; ok {
		ch.Close()
		delete(s.caches, hash)
	}
}

func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.caches {
		ch.Close()
	}
	return nil
}

func (s *Storage) CleanCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.caches {
		go ch.Clean()
	}
}

func (s *Storage) GetCache(hash string) *Cache {
	if ch, ok := s.caches[hash]; ok {
		return ch
	}
	return nil
}

func (s *Storage) GetStats() []CacheState {
	s.mu.Lock()
	defer s.mu.Unlock()
	cachesState := make([]CacheState, 0)
	for _, ch := range s.caches {
		cs := ch.GetState()
		cachesState = append(cachesState, cs)
	}
	sort.Slice(cachesState, func(i, j int) bool {
		return cachesState[i].Hash < cachesState[j].Hash
	})
	return cachesState
}
