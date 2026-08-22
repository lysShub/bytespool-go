package bytespool

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"unsafe"
)

func Test_SizeClasses(t *testing.T) {
	sizes, err := ProbeSizeClasses2()
	if err != nil {
		t.Fatal(err)
	}

	for i, e := range sizes[2:] {
		if e != classToSize[i] {
			t.Fatalf("class %d: size %d, want %d", i, classToSize[i], e)
		}
	}
}

func ProbeSizeClasses1() ([]int32, error) {
	path := filepath.Join(build.Default.GOROOT, "src", "internal", "runtime", "gc", "sizeclasses.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	const name = "SizeClassToSize"

	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, s := range gd.Specs {
			vs, ok := s.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || vs.Names[0].Name != name {
				continue
			}
			cl, ok := vs.Values[0].(*ast.CompositeLit)
			if !ok {
				return nil, fmt.Errorf("%s is not a CompositeLit", name)
			}
			vals := make([]int32, 0, len(cl.Elts))
			for _, e := range cl.Elts {
				bl, ok := e.(*ast.BasicLit)
				if !ok {
					return nil, fmt.Errorf("element of %s is not a BasicLit", name)
				}
				v, err := strconv.ParseInt(bl.Value, 10, 32)
				if err != nil {
					return nil, err
				}
				vals = append(vals, int32(v))
			}
			return vals, nil
		}
	}
	return nil, fmt.Errorf("%s not found", name)
}

func ProbeSizeClasses2() ([]int32, error) {
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
