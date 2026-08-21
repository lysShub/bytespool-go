package bytespool

// Pooler abstracts Get/Alloc/Put for a slice or object type.
type Pooler[T []e | *e, e any] interface {
	// Get returns pool-backed memory; not zeroes the memory.
	Get(len ...int) T
	// Alloc is like Get but zeroes the memory.
	Alloc(len ...int) T
	// Put returns memory to the pool.
	Put(v T)
}

type Pool[T []e | *e, e any] struct{}

var _ Pooler[[]byte, byte] = Pool[[]byte, byte]{}

func (Pool[T, e]) Get(len ...int) T   { return Get[T, e](len...) }
func (Pool[T, e]) Alloc(len ...int) T { return Alloc[T, e](len...) }
func (Pool[T, e]) Put(v T)            { Put[T, e](v) }
