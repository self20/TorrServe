package memcache

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sort"
	"sync"

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

	muPieces     sync.Mutex
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
	//c.bufferPool = NewBufferPool(int(c.pieceLength)) //c.capacity/int(info.PieceLength)+5
	c.hash = hash.HexString()

	for i := 0; i < c.pieceCount; i++ {
		c.pieces[info.Piece(i).Hash().HexString()] = &Piece{
			Id:     i,
			Length: info.Piece(i).Length(),
			Hash:   info.Piece(i).Hash().HexString(),
			//bufferPool: c.bufferPool,
		}
	}
}

func (c *Cache) Piece(m metainfo.Piece) storage.PieceImpl {
	if m.Index() >= len(c.pieces) {
		return nil
	}

	c.muPieces.Lock()
	defer c.muPieces.Unlock()
	if val, ok := c.pieces[m.Hash().HexString()]; ok {
		return val
	}
	return nil
}

func (c *Cache) Close() error {
	fmt.Println("Close cache for:", c.hash)
	c.Clean()
	return nil
}

func (c *Cache) Clean() {
	c.muPieces.Lock()
	defer c.muPieces.Unlock()
	for key, val := range c.pieces {
		if len(val.buffer) > 0 && val.complete {
			c.removePiece(key)
		}
	}
	releaseMemory()
}

func (c *Cache) GetState() CacheState {
	c.muPieces.Lock()
	defer c.muPieces.Unlock()

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
	curr := c.currentPiece
	end := c.endPiece
	sort.Slice(stats, func(i, j int) bool {
		id1 := stats[i].Id
		id2 := stats[j].Id

		id1in := id1 >= curr && id1 <= end
		id2in := id2 >= curr && id2 <= end
		if id1in && !id2in {
			return false
		}
		if !id1in && id2in {
			return true
		}

		if id1 > curr && id2 > curr {
			return id1 > id2
		}
		return id1 < id2
	})

	cState.Pieces = c.getRemoveItems()
	cState.Filled = c.filled
	return cState
}

func (c *Cache) CurrentRead(piece int) {
	c.muPieces.Lock()
	defer c.muPieces.Unlock()

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

	go c.cleanPieces()
}

func (c *Cache) cleanPieces() {
	c.muRemove.Lock()
	if c.isRemove {
		c.muRemove.Unlock()
		return
	}
	c.isRemove = true
	defer func() {
		c.isRemove = false
	}()
	c.muRemove.Unlock()

	c.muPieces.Lock()
	defer c.muPieces.Unlock()

	if c.capacity > 0 {
		removes := c.getRemoveItems()
		pos := 0
		for c.getFilled() > c.capacity && len(removes) > 0 && pos < len(removes) {
			c.removePiece(removes[pos].Hash)
			pos++
		}
		releaseMemory()
	}
}

func (c *Cache) removePiece(hash string) {
	if piece, ok := c.pieces[hash]; ok {
		piece.Release()
		st := fmt.Sprintf("%v\t%s\t%s\t%v", piece.Id, piece.accessed.Format("15:04:05.000"), piece.Hash, c.currentPiece)
		fmt.Println("Remove cache piece:", st)
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
