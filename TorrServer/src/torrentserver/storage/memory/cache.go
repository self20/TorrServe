package memory

import (
	"math"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"torrentserver/storage/state"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	"github.com/RoaringBitmap/roaring"
)

// Cache ...
type Cache struct {
	mu *sync.Mutex

	id        string
	running   bool
	capacity  int64
	filled    int64
	readahead int64

	policy Policy

	pieceCount    int
	pieceLength   int64
	piecePriority []int
	pieces        map[key]*Piece
	items         map[key]ItemState

	closing chan struct{}

	bufferSize int
	buffers    [][]byte
	positions  []*BufferPosition

	currentPiece int
	endPiece     int
}

// BufferPosition ...
type BufferPosition struct {
	Used  bool
	Index int
	Key   key
}

// CacheInfo is a container for basic active Cache into
type CacheInfo struct {
	Capacity int64
	Filled   int64
	Items    int
}

// ItemState ...
type ItemState struct {
	Accessed time.Time
	Size     int64
}

func NewCache(capacity int64, hash string, info *metainfo.Info) *Cache {
	ch := &Cache{
		capacity: capacity,
		id:       hash,
		mu:       &sync.Mutex{},
	}
	ch.Init(info)
	return ch
}

func (c *Cache) CurrentRead(piece int) {
	c.currentPiece = piece
	c.endPiece = piece + int(c.capacity/c.pieceLength)
	if c.endPiece > c.pieceCount {
		c.endPiece = c.pieceCount
	}
}

// SetCapacity ...
func (c *Cache) SetCapacity(capacity int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.capacity = capacity
}

// Piece ...
func (c *Cache) Piece(m metainfo.Piece) storage.PieceImpl {
	c.mu.Lock()
	defer c.mu.Unlock()

	if m.Index() >= len(c.pieces) {
		return nil
	}

	return c.pieces[key(m.Index())]
}

// Init creates buffers and underlying maps
func (c *Cache) Init(info *metainfo.Info) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[key]ItemState)
	c.policy = new(lru)

	c.pieceCount = info.NumPieces()
	c.pieceLength = info.PieceLength
	c.piecePriority = make([]int, c.pieceCount)

	// Using max possible buffers + 2
	c.bufferSize = int(math.Ceil(float64(c.capacity)/float64(c.pieceLength)) + 2)
	if c.bufferSize > c.pieceCount {
		c.bufferSize = c.pieceCount
	}
	c.readahead = int64(float64(c.capacity) * 0.33)

	c.buffers = make([][]byte, c.bufferSize)
	c.positions = make([]*BufferPosition, c.bufferSize)
	c.pieces = map[key]*Piece{}

	for i := 0; i < c.pieceCount; i++ {
		c.pieces[key(i)] = &Piece{
			c:        c,
			mu:       &sync.Mutex{},
			Position: -1,
			Index:    i,
			Key:      key(i),
			Length:   info.Piece(i).Length(),
			Hash:     info.Piece(i).Hash().HexString(),
			Chunks:   roaring.NewBitmap(),
		}
	}

	for i := range c.buffers {
		c.buffers[i] = make([]byte, c.pieceLength)
		c.positions[i] = &BufferPosition{}
	}
}

// Info returns information for Cache
func (c *Cache) Info() (ret CacheInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ret.Capacity = c.capacity
	ret.Filled = c.filled
	ret.Items = len(c.items)
	return
}

// Close ...
func (c *Cache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return nil
	}

	c.running = false
	c.Stop()
	return nil
}

// RemovePiece ...
func (c *Cache) RemovePiece(idx int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	k := key(idx)
	if _, ok := c.pieces[k]; ok {
		c.remove(k)
	}
}

func (c *Cache) GetState() state.CacheState {
	c.mu.Lock()
	defer c.mu.Unlock()

	cState := state.CacheState{}
	cState.Capacity = c.capacity
	cState.PiecesLength = c.pieceLength
	cState.PiecesCount = c.pieceCount
	cState.CurrentRead = c.currentPiece
	cState.EndRead = c.endPiece
	cState.Hash = c.id
	cState.Filled = c.filled

	stats := make([]state.ItemState, 0)
	for _, value := range c.pieces {
		ist := state.ItemState{
			Id:         value.Index,
			Hash:       value.Hash,
			Completed:  value.Completed,
			BufferSize: c.items[value.Key].Size,
			Accessed:   c.items[value.Key].Accessed,
		}
		stat := ist

		if stat.BufferSize > 0 {
			stats = append(stats, stat)
		}
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Accessed.Before(stats[j].Accessed)
	})
	cState.Pieces = stats
	cState.Filled = c.filled
	return cState
}

// Start is watching Cache statistics
func (c *Cache) Start() {
	c.running = true
	c.closing = make(chan struct{}, 1)
	progressTicker := time.NewTicker(1 * time.Second)

	defer progressTicker.Stop()
	defer close(c.closing)

	// var lastFilled int64

	for {
		select {
		case <-progressTicker.C:

		case <-c.closing:
			c.running = false
			return

		}
	}
}

// Stop ends progress timers, removes buffers, free memory to OS
func (c *Cache) Stop() {
	c.closing <- struct{}{}

	go func() {
		delay := time.NewTicker(1 * time.Second)
		defer delay.Stop()

		for {
			select {
			case <-delay.C:
				c.buffers = nil
				c.pieces = nil
				c.positions = nil

				debug.FreeOSMemory()

				return
			}
		}
	}()
}

func (c *Cache) remove(pi key) {
	// Don't allow to delete first piece, it's used everywhere
	if pi == 0 {
		return
	}

	if c.pieces[pi].Position != -1 {
		c.positions[c.pieces[pi].Position].Used = false
		c.positions[c.pieces[pi].Position].Index = 0
	}

	c.pieces[pi].Position = -1
	c.pieces[pi].Completed = false
	c.pieces[pi].Active = false
	c.pieces[pi].mu.Lock()
	c.pieces[pi].Reset()
	c.pieces[pi].mu.Unlock()

	c.updateItem(c.pieces[pi].Key, func(*ItemState, bool) bool {
		return false
	})
}

func (c *Cache) updateItem(k key, u func(*ItemState, bool) bool) {
	ii, ok := c.items[k]
	c.filled -= ii.Size
	if u(&ii, ok) {
		c.filled += ii.Size
		if int(k) != 0 {
			c.policy.Used(k, ii.Accessed)
		}
		c.items[k] = ii
	} else {
		c.policy.Forget(k)
		delete(c.items, k)
	}
	c.trimToCapacity()
}

// TrimToCapacity ...
func (c *Cache) TrimToCapacity() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.trimToCapacity()
}

func (c *Cache) trimToCapacity() {
	if c.capacity < 0 {
		return
	}
	for len(c.items) >= c.bufferSize {
		c.remove(c.policy.Choose().(key))
	}
}
