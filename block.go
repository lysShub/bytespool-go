package bytespool

import (
	"sync"
	"unsafe"
)

// Block is a pool-backed buffer that auto-grows and auto-shrinks.
type Block[T any] struct {
	_   [0]sync.Mutex // noCopy
	ptr unsafe.Pointer
	len int
	cap int
	cnt int // shrink counter
}

// Len returns the stored elements length.
func (b *Block[T]) Len() int { return b.len }

// Cap returns the buffer elements capacity.
func (b *Block[T]) Cap() int { return b.cap }

// Rem returns the remaining elements capacity (cap-len).
func (b *Block[T]) Rem() int { return b.cap - b.len }

// Release returns memory to the pool and resets b.
func (b *Block[T]) Release() {
	if b == nil || b.ptr == nil {
		return
	}
	put(b.ptr)
	*b = Block[T]{}
}

// Put stores a copy of s, auto-resizing as needed.
func (b *Block[T]) Put(s []T) *Block[T] {
	copy(b.Set(len(s)).Raw(), s)
	return b
}

// Set sets the length to n; if reserve is set, keep the original data when reallocating.
func (b *Block[T]) Set(n int, reserve ...bool) *Block[T] {
	if n < 0 {
		panic(n)
	}
	if n > b.cap {
		b.alloc(n, reserve...)
		return b
	}
	if n*2 <= b.cap {
		b.cnt++
	} else {
		b.cnt = 0
	}

	if b.cnt > 128 {
		b.alloc(n, reserve...)
	} else {
		b.len = n
	}
	return b
}
func (b *Block[T]) alloc(n int, reserved ...bool) {
	if b.ptr != nil {
		defer put(b.ptr)
	}
	var w = int(unsafe.Sizeof(*new(T)))

	var old []T
	if len(reserved) > 0 && reserved[0] {
		old = b.Raw()
	}

	b.len, b.cnt = n, 0
	b.ptr, b.cap = get(n*w, false)
	b.cap = b.cap / w
	if len(old) > 0 {
		copy(b.Raw(), old)
	}
}

// Raw returns the underlying slice.
func (b *Block[T]) Raw() []T {
	if b.ptr == nil {
		return nil
	}
	return unsafe.Slice((*T)(b.ptr), b.cap)[:b.len]
}
