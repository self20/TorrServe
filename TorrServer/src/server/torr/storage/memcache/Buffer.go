package memcache

import (
	"fmt"

	"server/utils"
)

type BufferPool struct {
	c       chan []byte
	bufSize int64
}

func NewBufferPool(bufSize int64, capacity int64) *BufferPool {
	bp := new(BufferPool)
	bp.bufSize = bufSize
	bp.c = make(chan []byte, int(capacity/bufSize)+3)
	for i := 0; i < int(capacity/bufSize)+3; i++ {
		bp.PutBuffer(make([]byte, bufSize))
	}
	return bp
}

func (b *BufferPool) GetBuffer() (buf []byte) {
	select {
	case buf = <-b.c:
	default:
		buf = make([]byte, b.bufSize)
		fmt.Println("Create buffer", len(b.c))
	}
	return
}

func (b *BufferPool) PutBuffer(buf []byte) {
	select {
	case b.c <- buf:
		fmt.Println("Save buffer", len(b.c))
	default:
		fmt.Println("Slow delete memory")
		buf = nil
		utils.FreeOSMemGC()
	}
}

func (b *BufferPool) Len() int {
	return len(b.c)
}
