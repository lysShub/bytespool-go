package bytespool

import (
	"unsafe"
)

// Get returns a pool-backed slice or object of type T. The memory is not zeroed.
func Get[T []e | *e, e any](len ...int) T {
	return fetch[T, e](false, len...)
}

// Alloc is like [Get] but zeroes the memory.
func Alloc[T []e | *e, e any](len ...int) T {
	return fetch[T, e](true, len...)
}

// Put returns a slice or object obtained from [Get]/[Alloc] to the pool.
func Put[T []e | *e, e any](v T) {
	put(*(*unsafe.Pointer)(unsafe.Pointer(&v)))
}

const (
	hdrsize        = 8
	magic   uint64 = 0x00fe1b79aec72663
)

func fetch[T []e | *e, e any](clr bool, n ...int) T {
	var w = int(unsafe.Sizeof(*new(e)))

	if unsafe.Sizeof(*new(T)) == unsafe.Sizeof(uintptr(0)) {
		// object
		if Debug && len(n) > 0 {
			panic(n[0])
		}
		p, _ := get(w, clr)
		return *(*T)(unsafe.Pointer(&p))
	} else {
		// slice
		var bytes = 0
		if len(n) > 0 {
			bytes = w * n[0]
		}
		p, c := get(bytes, clr)
		s := unsafe.Slice((*e)(p), c/w)

		if len(n) > 0 {
			s = s[:n[0]]
		}
		return *(*T)(unsafe.Pointer(&s))
	}
}

func enc(p unsafe.Pointer, poolIdx poolIdx) {
	*(*uint64)(p) = uint64(poolIdx)<<56 + magic
}
func dec(p unsafe.Pointer) (ok bool, idx poolIdx) {
	v := *(*uint64)(p)
	idx = poolIdx(v >> 56)
	return ((v<<8)>>8) == magic && idx <= poolIdxExceed, idx
}

func get(bytes int, clr bool) (p unsafe.Pointer, c int) {
	bytes += hdrsize

	idx := newPoolIdx(bytes)
	if idx > maxPoolIdx {
		if Debug {
			debug_log_get_exceed(bytes)
		}
		// block larger than the largest pool, allocate directly
		b := make([]byte, bytes)
		p = unsafe.Pointer(unsafe.SliceData(b))
		c = cap(b)

		enc(p, poolIdxExceed)
	} else {
		p = pools[idx].Get().(unsafe.Pointer)
		c = idx.bytes()

		enc(p, idx)
	}
	if Debug {
		debug_get(uintptr(p), idx)
	}

	p = unsafe.Add(p, hdrsize)
	c -= hdrsize
	if clr {
		clear(unsafe.Slice((*byte)(p), c))
	}
	return p, c
}

func put(ptr unsafe.Pointer) {
	if ptr == nil {
		if Debug {
			debug_log_put_nil()
		}
		return
	}

	ptr = unsafe.Add(ptr, -hdrsize)
	ok, i := dec(ptr)
	if !ok {
		if Debug {
			debug_log_put_invalid(uintptr(ptr))
		}
	} else {
		if Debug {
			debug_put(uintptr(ptr), i)
		}
		if i <= maxPoolIdx {
			pools[i].Put(ptr)
		}
	}
}
