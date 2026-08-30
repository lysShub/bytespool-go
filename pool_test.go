package bytespool

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"runtime"
	"slices"
	"testing"
	"unsafe"
)

type Object struct{ A, B int16 }

const maxBlocksize = large4096k

func Test_Get(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		s := Get[[]byte, byte](1024)
		if len(s) != 1024 {
			t.Fatalf("len = %d, want 1024", len(s))
		}
		if cap(s) < 1024 {
			t.Fatalf("cap = %d, want >= 1024", cap(s))
		}
		Put[[]byte, byte](s)
	})

	t.Run("object", func(t *testing.T) {
		p := Get[*Object, Object]()
		if p == nil {
			t.Fatal("got nil")
		}
		Put[*Object, Object](p)
	})

	t.Run("reuse", func(t *testing.T) {
		s1 := Get[[]byte, byte](111)
		rand.Read(s1)
		bak := slices.Clone(s1)
		Put[[]byte, byte](s1)

		s2 := Get[[]byte, byte](111)
		if !bytes.Equal(bak, s2) {
			t.Fatal("reused buffer content mismatch")
		}
	})

	t.Run("exceed", func(t *testing.T) {
		s := Get[[]byte, byte](maxBlocksize + 1)
		if len(s) < maxBlocksize+1 {
			t.Fatalf("len = %d, want >= %d", len(s), maxBlocksize+1)
		}
		Put[[]byte, byte](s)
	})
}

func Test_Alloc(t *testing.T) {
	t.Run("slice/zeroed", func(t *testing.T) {
		s := Alloc[[]byte, byte](1024)
		if len(s) < 1024 {
			t.Fatalf("len = %d, want >= 1024", len(s))
		}
		for _, b := range s {
			if b != 0 {
				t.Fatalf("byte = %d, want 0", b)
			}
		}
		Put[[]byte, byte](s)
	})

	t.Run("obj/zeroed", func(t *testing.T) {
		p := Alloc[*Object, Object]()
		if p == nil {
			t.Fatal("got nil")
		}
		if *p != (Object{}) {
			t.Fatalf("got %+v, want zero value", *p)
		}
		Put[*Object, Object](p)
	})
}

func Test_Put(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		Put[[]byte, byte]([]byte(nil))
		Put[*int, int](nil)
	})
}

func Test_newPoolIdx(t *testing.T) {
	probes := []int{hdrsize}
	for _, sz := range classToSize {
		probes = append(probes, int(sz), int(sz)+1)
	}
	for bytes := hdrsize; bytes <= maxBlocksize; bytes += 997 {
		probes = append(probes, bytes)
	}
	probes = append(probes, maxBlocksize, maxBlocksize+1)

	for _, bytes := range probes {
		idx := newPoolIdx(bytes)
		if bytes > maxBlocksize {
			if idx != poolIdxExceed {
				t.Fatalf("idx = %d, want %d", idx, poolIdxExceed)
			}
			continue
		}
		// the returned pool must fit bytes and be the smallest that does
		if idx > maxPoolIdx {
			t.Fatalf("idx = %d, want <= %d", idx, maxPoolIdx)
		}
		if int(classToSize[idx]) < bytes {
			t.Fatalf("class size %d < requested %d", classToSize[idx], bytes)
		}
		if idx > 0 {
			if int(classToSize[idx-1]) >= bytes {
				t.Fatalf("previous class size %d >= requested %d", classToSize[idx-1], bytes)
			}
		}
	}
}

func Test_Alignment(t *testing.T) {
	for _, sz := range classToSize {
		n := int(sz) - hdrsize
		if n <= 0 {
			continue
		}
		s := Get[[]byte, byte](n)
		p := unsafe.Pointer(unsafe.SliceData(s))
		if uintptr(p)%8 != 0 {
			t.Fatalf("pool size %d: block pointer %p not 8-aligned", sz, p)
		}
		Put[[]byte, byte](s)
	}
}

func Benchmark_Pool(b *testing.B) {
	var r1, r2 float64
	var procs = runtime.GOMAXPROCS(0)

	b.Run("serialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := Get[[]byte, byte](111)
			Put[[]byte, byte](s)
		}
		r1 = float64(b.Elapsed()) / float64(b.N)
	})
	b.Run(fmt.Sprintf("parallel___%d", procs), func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				s := Get[[]byte, byte](111)
				Put[[]byte, byte](s)
			}
		})
		r2 = float64(b.Elapsed()) / float64(b.N)
	})

	// validate global singleton [sync.Pool] would not be bottleneck
	if r2 > r1*1.5 {
		b.Fatal("parallel bottleneck")
	}
}
