package memcache

import (
	"fmt"
	"sort"
	"sync"

	"torrentserver/storage/state"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

type Cache struct {
	storage.TorrentImpl

	s *Storage

	capacity int64
	filled   int64
	hash     string

	pieceLength int64
	pieceCount  int

	muPiece  sync.Mutex
	muRemove sync.Mutex
	isRemove bool

	pieces       map[string]*Piece
	currentPiece int
	endPiece     int
}

func NewCache(capacity int64, storage *Storage) *Cache {
	ret := &Cache{
		capacity: capacity,
		filled:   0,
		pieces:   make(map[string]*Piece),
		s:        storage,
	}

	return ret
}

func (c *Cache) Init(info *metainfo.Info, hash metainfo.Hash) {
	fmt.Println("Create cache for:", info.Name)
	//Min capacity of 2 pieces length
	cap := info.PieceLength * 2
	if c.capacity < cap {
		c.capacity = cap
	}
	c.pieceLength = info.PieceLength
	c.pieceCount = info.NumPieces()
	c.hash = hash.HexString()

	for i := 0; i < c.pieceCount; i++ {
		c.pieces[info.Piece(i).Hash().HexString()] = &Piece{
			Id:     i,
			Length: info.Piece(i).Length(),
			Hash:   info.Piece(i).Hash().HexString(),
			cache:  c,
		}
	}
}

func (c *Cache) Piece(m metainfo.Piece) storage.PieceImpl {
	if m.Index() >= len(c.pieces) {
		return nil
	}

	c.muPiece.Lock()
	defer c.muPiece.Unlock()
	if val, ok := c.pieces[m.Hash().HexString()]; ok {
		return val
	}
	return nil
}

func (c *Cache) Close() error {
	c.isRemove = false
	fmt.Println("Close cache for:", c.hash)
	if _, ok := c.s.caches[c.hash]; ok {
		delete(c.s.caches, c.hash)
	}

	c.pieces = nil
	releaseMemory()
	return nil
}

func (c *Cache) Clean() {
	c.pieces = make(map[string]*Piece)
	releaseMemory()
}

func (c *Cache) GetState() state.CacheState {
	cState := state.CacheState{}
	cState.Capacity = c.capacity
	cState.PiecesLength = c.pieceLength
	cState.PiecesCount = c.pieceCount
	cState.CurrentRead = c.currentPiece
	cState.EndRead = c.endPiece
	cState.Hash = c.hash
	cState.Filled = c.filled

	stats := make([]state.ItemState, 0)
	c.muPiece.Lock()
	for _, value := range c.pieces {
		stat := value.Stat()
		if stat.BufferSize > 0 {
			stats = append(stats, stat)
		}
	}
	c.muPiece.Unlock()
	sort.Slice(stats, func(i, j int) bool {
		id1 := stats[i].Id
		id2 := stats[j].Id
		return id1 < id2
	})
	cState.Pieces = stats
	return cState
}

func (c *Cache) CurrentRead(piece int) {
	c.currentPiece = piece
	c.endPiece = piece + int(c.capacity/c.pieceLength)
	if c.endPiece > c.pieceCount {
		c.endPiece = c.pieceCount
	}
}

func (c *Cache) cleanPieces() {
	if c.isRemove {
		return
	}
	c.muRemove.Lock()
	if c.isRemove {
		c.muRemove.Unlock()
		return
	}
	c.isRemove = true
	defer func() { c.isRemove = false }()
	c.muRemove.Unlock()

	remPieces := c.getRemPieces()
	if len(remPieces) > 0 && c.capacity < c.filled {
		remCount := int((c.filled - c.capacity) / c.pieceLength)
		if remCount < 1 {
			remCount = 1
		}
		if remCount > len(remPieces) {
			remCount = len(remPieces)
		}

		remPieces = remPieces[:remCount]

		for _, p := range remPieces {
			c.removePiece(p)
		}
	}
}

func (c *Cache) getRemPieces() []*Piece {
	if c.currentPiece == 0 && c.endPiece == 0 {
		return nil
	}

	pieces := make([]*Piece, 0)
	var curr *Piece
	fill := int64(0)
	for _, v := range c.pieces {
		if v.Size > 0 {
			pieces = append(pieces, v)
			fill += v.Size
		}
		if v.Id == c.currentPiece {
			curr = v
		}
	}
	c.filled = fill
	if curr == nil {
		return nil
	}
	sort.Slice(pieces, func(i, j int) bool {
		return pieces[i].accessed.Before(pieces[j].accessed)
	})
	pos := 0
	for i, v := range pieces {
		if v.accessed.UnixNano() >= curr.accessed.UnixNano() {
			pos = i
			break
		}
	}
	return pieces[:pos]
}

func (c *Cache) removePiece(piece *Piece) {
	c.muPiece.Lock()
	defer c.muPiece.Unlock()
	piece.Release()
	st := fmt.Sprintf("%v\t%s\t%s\t%v", piece.Id, piece.accessed.Format("15:04:05.000"), piece.Hash, c.currentPiece)
	fmt.Println("Remove cache piece:", st)
	releaseMemory()
}
