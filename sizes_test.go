package bytespool

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"testing"
)

func Test_SizeClasses(t *testing.T) {
	path := filepath.Join(build.Default.GOROOT, "src", "internal", "runtime", "gc", "sizeclasses.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	extract := func(name string) []uint64 {
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
					t.Fatalf("%s is not a CompositeLit", name)
				}
				vals := make([]uint64, 0, len(cl.Elts))
				for _, e := range cl.Elts {
					bl, ok := e.(*ast.BasicLit)
					if !ok {
						t.Fatalf("element of %s is not a BasicLit", name)
					}
					v, err := strconv.ParseUint(bl.Value, 10, 64)
					if err != nil {
						t.Fatal(err)
					}
					vals = append(vals, v)
				}
				return vals
			}
		}
		t.Fatalf("%s not found", name)
		return nil
	}

	scls := extract("SizeClassToSize")
	for i, e := range scls[2:] {
		if e != uint64(classToSize[i]) {
			t.Fatalf("class %d: size %d, want %d", i, classToSize[i], e)
		}
	}
}
