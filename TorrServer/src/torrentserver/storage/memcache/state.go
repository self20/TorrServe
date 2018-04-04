package memcache

import (
	"time"
)

type CacheState struct {
	Hash         string
	Capacity     int
	Filled       int
	PiecesLength int
	PiecesCount  int
	CurrentRead  int
	Pieces       []ItemState
}

type ItemState struct {
	Id         int
	Accessed   time.Time
	BufferSize int
	Completed  bool
	Hash       string
}

func (c CacheState) FindItemId(id int) *ItemState {
	for _, itm := range c.Pieces {
		if itm.Id == id {
			return &itm
		}
	}
	return nil
}
