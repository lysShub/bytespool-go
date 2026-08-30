package bytespool_test

import (
	"fmt"

	"github.com/lysShub/bytespool-go"
)

func ExampleGet_slice() {
	s := bytespool.Get[[]byte, byte](8)
	copy(s, "hello")
	fmt.Println(len(s), string(s[:5]))
	bytespool.Put[[]byte, byte](s)

	// Output: 8 hello
}

func ExampleGet_object() {
	type Frame struct {
		Type uint16
		Len  uint32
		Body [1024]byte
	}

	p := bytespool.Get[*Frame, Frame]()
	defer bytespool.Put[*Frame, Frame](p)

	p.Type = 1
	p.Len = 3
	copy(p.Body[:], "abc")
	fmt.Println(p.Type, p.Len, string(p.Body[:p.Len]))

	// Output: 1 3 abc
}
