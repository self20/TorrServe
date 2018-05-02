package memcache

import (
	"runtime"
	"runtime/debug"
)

func releaseMemory() {
	runtime.GC()
	debug.FreeOSMemory()
}
