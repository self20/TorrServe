package memcache

import "torrentserver/utils"

type BufferPool struct {
	c       chan []byte
	bufSize int64
}

func NewBufferPool(bufSize int64, capacity int64) *BufferPool {
	bp := new(BufferPool)
	bp.bufSize = bufSize
	bp.c = make(chan []byte, int(capacity/bufSize)+3)
	return bp
}

func (b *BufferPool) GetBuffer() (buf []byte) {
	select {
	case buf = <-b.c:
	default:
		buf = make([]byte, b.bufSize)
	}
	return
}

func (b *BufferPool) PutBuffer(buf []byte) {
	select {
	case b.c <- buf:
	default:
		utils.ReleaseMemory()
	}
}
