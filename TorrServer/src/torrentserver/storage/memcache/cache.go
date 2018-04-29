package memcache

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

type Cache struct {
	storage.TorrentImpl

	capacity int
	filled   int
	hash     string

	pieceLength int64
	pieceCount  int

	muRemove sync.Mutex
	isRemove bool

	pieces       map[string]*Piece
	currentPiece int
	endPiece     int
}

func NewCache(capacity int) *Cache {
	ret := &Cache{
		capacity: capacity,
		filled:   0,
		pieces:   make(map[string]*Piece),
	}

	return ret
}

func (c *Cache) Init(info *metainfo.Info, hash metainfo.Hash) {
	fmt.Println("Create cache for:", info.Name)
	//Min capacity of 10 pieces length
	cap := int(info.PieceLength * 10)
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
		}
	}
	go c.cleanPieces()
}

func (c *Cache) Piece(m metainfo.Piece) storage.PieceImpl {
	if m.Index() >= len(c.pieces) {
		return nil
	}

	if val, ok := c.pieces[m.Hash().HexString()]; ok {
		return val
	}
	return nil
}

func (c *Cache) Close() error {
	c.isRemove = false
	fmt.Println("Close cache for:", c.hash)
	c.Clean()
	return nil
}

func (c *Cache) Clean() {
	for key, val := range c.pieces {
		if len(val.buffer) > 0 {
			c.removePiece(key)
		}
	}
}

func (c *Cache) GetState() CacheState {
	cState := CacheState{}
	cState.Capacity = c.capacity
	cState.PiecesLength = int(c.pieceLength)
	cState.PiecesCount = c.pieceCount
	cState.CurrentRead = c.currentPiece
	cState.Hash = c.hash

	stats := make([]ItemState, 0)
	c.filled = 0
	for _, value := range c.pieces {
		stat := value.Stat()
		c.filled += stat.BufferSize
		if stat.BufferSize > 0 {
			stats = append(stats, stat)
		}
	}
	sort.Slice(stats, func(i, j int) bool {
		id1 := stats[i].Id
		id2 := stats[j].Id
		return id1 < id2
	})
	cState.PiecesInCache = stats
	cState.PiecesForDel = c.getRemoveItems()
	cState.Filled = c.filled
	return cState
}

func (c *Cache) CurrentRead(piece int) {
	c.currentPiece = piece
	c.endPiece = piece + (c.capacity / int(c.pieceLength))
	if c.endPiece > c.pieceCount {
		c.endPiece = c.pieceCount
	}

	c.filled = 0
	for _, pi := range c.pieces {
		stat := pi.Stat()
		c.filled += stat.BufferSize
	}

	if c.filled < c.capacity {
		return
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

	for c.isRemove {
		if c.capacity > 0 {
			removes := c.getRemoveItems()
			pos := 0
			for c.getFilled() > c.capacity && len(removes) > 0 && pos < len(removes) {
				c.removePiece(removes[pos].Hash)
				pos++
			}
		}
		time.Sleep(time.Second)
	}
}

func (c *Cache) removePiece(hash string) {
	if piece, ok := c.pieces[hash]; ok {
		piece.Release()
		st := fmt.Sprintf("%v\t%s\t%s\t%v", piece.Id, piece.accessed.Format("15:04:05.000"), piece.Hash, c.currentPiece)
		fmt.Println("Remove cache piece:", st)
		releaseMemory()
	}
}

func (c *Cache) getRemoveItems() []ItemState {
	removes := make([]ItemState, 0)
	c.filled = 0
	for _, pi := range c.pieces {
		stat := pi.Stat()
		c.filled += stat.BufferSize

		if (pi.Id < c.currentPiece || pi.Id > c.endPiece) && pi.Id > 0 && len(pi.buffer) > 0 {
			removes = append(removes, stat)
		}
	}
	curr := c.currentPiece
	sort.Slice(removes, func(i, j int) bool {
		id1 := removes[i].Id
		id2 := removes[j].Id

		if id1 > curr && id2 > curr {
			return id1 > id2
		}
		return id1 < id2
	})
	return removes
}

func (c *Cache) getFilled() int {
	c.filled = 0
	for _, pi := range c.pieces {
		stat := pi.Stat()
		c.filled += stat.BufferSize
	}
	return c.filled
}

func releaseMemory() {
	runtime.GC()
	debug.FreeOSMemory()
}
