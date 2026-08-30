# bytespool-go

bytespool is a generic memory pool built on top of `sync.Pool`. It reuses `[]T` slices and `*T` objects obtained by [`Get`]/[`Alloc`] and returned by [`Put`], reducing GC pressure and allocations in hot paths.

## Features

- **generic**: works with both `[]T` slices and `*T` objects
- **easy `Put`**: pass only the head pointer of the memory to `Put`
- **concurrency**: the pool is safe for concurrent use
- **aligned**: block sizes are aligned with Go's runtime size classes
- **debug mode**: build with `-tags debug` to detect misuse (via slog.Warn):
  1. double `Put`
  2. invalid `Put`
  3. `Put` after GC

## Example

Get a `[]byte` slice:

```go
s := bytespool.Get[[]byte, byte](16)
defer bytespool.Put[[]byte, byte](s)

copy(s, "hello")
```

Get a `*T` object. Only pool variables that escape to the heap — for small values that stay on the stack, pooling is a pessimization:

```go
type Frame struct {
	Type uint16
	Len  uint32
	Body [1024]byte
}

p := bytespool.Get[*Frame, Frame]()
defer bytespool.Put[*Frame, Frame](p)

p.Type = 1
copy(p.Body[:], "hello")
```

Use [`Alloc`] instead of [`Get`] when zeroed memory is required.
