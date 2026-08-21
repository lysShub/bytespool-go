package bytespool

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func Test_Block(t *testing.T) {
	t.Run("buff/release", func(t *testing.T) {
		p := Block[byte]{}
		defer p.Release()

		var b = make([]byte, 123)
		rand.Read(b)
		p.Put(b)
		if !bytes.Equal(b, p.Raw()) {
			t.Fatal("Raw() mismatch")
		}

		p.Release()
		if p.ptr != nil {
			t.Fatal("ptr not nil after Release")
		}
	})

	t.Run("alloc/large", func(t *testing.T) {
		p := Block[byte]{}
		defer p.Release()

		var b1 = make([]byte, 123)
		rand.Read(b1)
		p.Put(b1)
		ptr1 := p.ptr

		var b2 = make([]byte, 1396)
		p.Put(b2)
		if ptr1 == p.ptr {
			t.Fatal("ptr not reallocated for larger buffer")
		}
	})

	t.Run("alloc/small", func(t *testing.T) {
		p := Block[byte]{}
		defer p.Release()

		p.Put(make([]byte, 1024))
		if p.cap != 1152-hdrsize {
			t.Fatalf("cap = %d, want %d", p.cap, 1152-hdrsize)
		}
		prt1 := p.ptr

		for i := 0; i < 129; i++ {
			p.Put(make([]byte, 111))
			if i == 0 {
				if prt1 != p.ptr {
					t.Fatal("ptr changed on first small buff")
				}
			}
		}
		if p.cap != 128-hdrsize {
			t.Fatalf("cap = %d, want %d", p.cap, 128-hdrsize)
		}
	})
}
