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

	t.Run("alloc/shrink-reserve", func(t *testing.T) {
		p := Block[byte]{}
		defer p.Release()

		var b = make([]byte, 1024)
		rand.Read(b)
		p.Put(b)
		ptr1 := p.ptr

		for i := 0; i < 129; i++ {
			p.Set(512, true)
		}
		if p.ptr == ptr1 {
			t.Fatal("ptr not reallocated on shrink")
		}
		if p.Len() != 512 {
			t.Fatalf("Len = %d, want 512", p.Len())
		}
		if !bytes.Equal(b[:512], p.Raw()) {
			t.Fatal("old prefix not kept on shrink")
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

	t.Run("append/grow", func(t *testing.T) {
		p := Block[byte]{}
		defer p.Release()

		p.Append(0)
		cap0 := p.cap
		old := append([]byte(nil), p.Raw()...)

		p.Append(1, 2, 3)
		if p.Len() != 4 {
			t.Fatalf("Len = %d, want 4", p.Len())
		}
		if p.cap != cap0 {
			t.Fatalf("cap changed on in-cap append: %d -> %d", cap0, p.cap)
		}
		for i, v := range old {
			if p.Raw()[i] != v {
				t.Fatalf("old data lost at %d: got %d, want %d", i, p.Raw()[i], v)
			}
		}
		for i, v := range []byte{1, 2, 3} {
			if p.Raw()[len(old)+i] != v {
				t.Fatalf("appended data wrong at %d: got %d, want %d", len(old)+i, p.Raw()[len(old)+i], v)
			}
		}
	})

	t.Run("append/overflow", func(t *testing.T) {
		p := Block[int]{}
		defer p.Release()

		var old = make([]int, 10)
		for i := range old {
			old[i] = i
		}
		p.Append(old...)

		var more = make([]int, 100000)
		for i := range more {
			more[i] = 10 + i
		}
		p.Append(more...)

		if p.Len() != len(old)+len(more) {
			t.Fatalf("Len = %d, want %d", p.Len(), len(old)+len(more))
		}
		for i, v := range old {
			if p.Raw()[i] != v {
				t.Fatalf("old data lost at %d: got %d, want %d", i, p.Raw()[i], v)
			}
		}
		for i, v := range more {
			if p.Raw()[len(old)+i] != v {
				t.Fatalf("appended data wrong at %d: got %d, want %d", len(old)+i, p.Raw()[len(old)+i], v)
			}
		}
	})
}
