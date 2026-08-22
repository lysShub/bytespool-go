package bytespool

import (
	"errors"
	"runtime"
	"testing"
	"unsafe"
)

func Test_SizeClasses(t *testing.T) {
	sizes, err := ProbeSizeClasses()
	if err != nil {
		t.Fatal(err)
	}

	for i, e := range sizes[2:] {
		if e != classToSize[i] {
			t.Fatalf("class %d: size %d, want %d", i, classToSize[i], e)
		}
	}
}

func ProbeSizeClasses() ([]int32, error) {
	var sizes = []int32{0, 8}
	for n := 16; n <= 32768; {
		s := probeAllocSize(n)
		if s < uintptr(n) {
			return nil, errors.New("probeAllocSize fail")
		}
		if s != uintptr(n) {
			n = int(s)
			continue
		}
		sizes = append(sizes, int32(n))
		n++
	}
	return sizes, nil
}
func probeAllocSize(n int) uintptr {
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	const samples = 64
	ptrs := make([]unsafe.Pointer, 0, samples)
	for i := 0; i < samples; i++ {
		b := make([]byte, n)
		ptrs = append(ptrs, unsafe.Pointer(unsafe.SliceData(b)))
	}
	runtime.KeepAlive(ptrs)
	runtime.ReadMemStats(&after)

	return uintptr((after.TotalAlloc - before.TotalAlloc) / samples)
}
