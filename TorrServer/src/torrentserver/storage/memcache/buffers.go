package memcache

import "sync"

//
//type BufferPool struct {
//	chanPool chan []byte
//	buffSize int
//}
//
//func NewBufferPool(poolSize int, bufSize int) *BufferPool {
//	bp := new(BufferPool)
//	bp.chanPool = make(chan []byte, poolSize)
//	bp.buffSize = bufSize
//	return bp
//}
//
//func (b *BufferPool) Get() []byte {
//	select {
//	case bt, ok := <-b.chanPool:
//		if ok {
//			if cap(bt) < b.buffSize {
//				return make([]byte, b.buffSize)
//			}
//			return bt[:b.buffSize]
//		}
//	default:
//	}
//	return make([]byte, size)
//}
//
//func (b *BufferPool) Put(buf []byte) {
//	buf = buf[:0]
//	select {
//	case b.chanPool <- buf:
//	default:
//	}
//	return
//}

type BufferPool struct {
	bytesPool sync.Pool
	buffSize  int
}

func NewBufferPool(bufSize int) *BufferPool {
	bp := new(BufferPool)
	bp.buffSize = bufSize
	return bp
}

func (b *BufferPool) Get() []byte {
	ifc := b.bytesPool.Get()
	if ifc != nil {
		buf := ifc.([]byte)
		if cap(buf) < b.buffSize {
			return make([]byte, b.buffSize)
		}
		return buf[:b.buffSize]
	}
	return make([]byte, b.buffSize)
}

func (b *BufferPool) Put(buf []byte) {
	b.bytesPool.Put(buf)
}
