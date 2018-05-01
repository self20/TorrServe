package state

import (
	"time"
)

type CacheState struct {
	Hash          string
	Capacity      int64
	Filled        int64
	PiecesLength  int64
	PiecesCount   int
	CurrentRead   int
	EndRead       int
	PiecesForDel  []ItemState
	PiecesInCache []ItemState
}

type ItemState struct {
	Id         int
	Accessed   time.Time
	BufferSize int64
	Completed  bool
	Hash       string
}

func (c CacheState) FindItemId(id int) *ItemState {
	for _, itm := range c.PiecesForDel {
		if itm.Id == id {
			return &itm
		}
	}
	return nil
}
