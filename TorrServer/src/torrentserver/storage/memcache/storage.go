package memcache

import (
	"sort"
	"sync"

	"torrentserver/settings"
	"torrentserver/storage/memory"
	"torrentserver/storage/state"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

type Storage struct {
	storage.ClientImpl

	caches   map[string]*Cache
	elcaches map[string]*memory.Cache
	capacity int64
	mu       sync.Mutex
}

func NewStorage(capacity int64) *Storage {
	stor := new(Storage)
	stor.capacity = capacity
	stor.caches = make(map[string]*Cache)
	stor.elcaches = make(map[string]*memory.Cache)
	return stor
}

func (s *Storage) OpenTorrent(info *metainfo.Info, infoHash metainfo.Hash) (storage.TorrentImpl, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if settings.Get().IsElementumCache {
		ch := memory.NewCache(int64(s.capacity), infoHash.HexString(), info)
		s.elcaches[infoHash.HexString()] = ch
		return ch, nil
	} else {
		ch := NewCache(s.capacity, s)
		ch.Init(info, infoHash)
		s.caches[infoHash.HexString()] = ch
		return ch, nil
	}
}

func (s *Storage) CloseByHash(hash string) {
	if settings.Get().IsElementumCache {
		if s.elcaches == nil {
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if ch, ok := s.elcaches[hash]; ok {
			ch.Close()
			delete(s.elcaches, hash)
		}
	} else {
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
}

func (s *Storage) Close() error {
	if settings.Get().IsElementumCache {
		s.mu.Lock()
		defer s.mu.Unlock()
		for _, ch := range s.elcaches {
			ch.Close()
		}
		return nil
	} else {
		s.mu.Lock()
		defer s.mu.Unlock()
		for _, ch := range s.caches {
			ch.Close()
		}
		return nil
	}
}

func (s *Storage) CleanCache() {
	if !settings.Get().IsElementumCache {
		s.mu.Lock()
		defer s.mu.Unlock()
		for _, ch := range s.caches {
			go ch.Clean()
		}
	}
}

func (s *Storage) GetCache(hash string) *Cache {
	if ch, ok := s.caches[hash]; ok {
		return ch
	}
	return nil
}

func (s *Storage) GetElCache(hash string) *memory.Cache {
	if ch, ok := s.elcaches[hash]; ok {
		return ch
	}
	return nil
}

func (s *Storage) GetStats() []state.CacheState {
	if settings.Get().IsElementumCache {
		s.mu.Lock()
		defer s.mu.Unlock()
		cachesState := make([]state.CacheState, 0)
		for _, ch := range s.elcaches {
			cs := ch.GetState()
			cachesState = append(cachesState, cs)
		}
		sort.Slice(cachesState, func(i, j int) bool {
			return cachesState[i].Hash < cachesState[j].Hash
		})
		return cachesState
	} else {
		s.mu.Lock()
		defer s.mu.Unlock()
		cachesState := make([]state.CacheState, 0)
		for _, ch := range s.caches {
			cs := ch.GetState()
			cachesState = append(cachesState, cs)
		}
		sort.Slice(cachesState, func(i, j int) bool {
			return cachesState[i].Hash < cachesState[j].Hash
		})
		return cachesState
	}
}
