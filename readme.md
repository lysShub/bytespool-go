# bytespool-go

bytespool is a generic memory pool built on top of `sync.Pool`. It reuses `[]T` slices and `*T` objects obtained by [`Get`]/[`Alloc`] and returned by [`Put`], reducing GC pressure and allocations in hot paths.

## Features

- **generic**: works with both `[]T` slices and `*T` objects
- **easy put**: pass only the head pointer of the memory to `Put`
- **concurrency**: the pool is safe for concurrent use
- **aligned**: block sizes are aligned with Go's runtime size classes
- **debug mode**: build with `-tags debug` to detect misuse. Warnings are reported via `slog.Warn`:
  1. double `Put`
  2. invalid `Put`
  3. `Put` after GC

## Doc

[https://pkg.go.dev/github.com/lysShub/bytespool-go](https://pkg.go.dev/github.com/lysShub/bytespool-go)
