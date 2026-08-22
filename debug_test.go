//go:build debug
// +build debug

package bytespool

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"unsafe"
)

func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

func Test_Debug_PutNil(t *testing.T) {
	DebugClear()
	buf := captureLog(t)

	Put[[]byte, byte]([]byte(nil))

	if !strings.Contains(buf.String(), "bytespool put nil") {
		t.Fatalf("log not printed: %s", buf.String())
	}
}

func Test_Debug_PutInvalid(t *testing.T) {
	DebugClear()
	buf := captureLog(t)

	Put[[]byte, byte](make([]byte, 16))

	if !strings.Contains(buf.String(), "bytespool put invalid") {
		t.Fatalf("log not printed: %s", buf.String())
	}
}

func Test_Debug_PutUntracked(t *testing.T) {
	DebugClear()
	buf := captureLog(t)

	s := Get[[]byte, byte](16)
	DebugClear()
	Put[[]byte, byte](s)

	if !strings.Contains(buf.String(), "bytespool put untracked") {
		t.Fatalf("log not printed: %s", buf.String())
	}
}

func Test_Debug_GetExceed(t *testing.T) {
	DebugClear()
	buf := captureLog(t)

	s := Get[[]byte, byte](maxBlocksize + 1)
	Put[[]byte, byte](s)

	if !strings.Contains(buf.String(), "bytespool get exceed") {
		t.Fatalf("log not printed: %s", buf.String())
	}
}

func Test_Debug_PutDamaged(t *testing.T) {
	DebugClear()
	buf := captureLog(t)

	s := Get[[]byte, byte](16)
	hdr := unsafe.Pointer(uintptr(unsafe.Pointer(unsafe.SliceData(s))) - hdrsize)
	old := *(*byte)(unsafe.Add(hdr, 7))
	*(*byte)(unsafe.Add(hdr, 7)) = old ^ 1
	Put[[]byte, byte](s)

	if !strings.Contains(buf.String(), "bytespool put damaged") {
		t.Fatalf("log not printed: %s", buf.String())
	}
}

func Test_Debug_PutAfterGC(t *testing.T) {
	DebugClear()
	buf := captureLog(t)

	p := unsafe.Pointer(new([16]byte))
	debug_get(uintptr(p), 0)
	debug_get(uintptr(p), 0)

	if !strings.Contains(buf.String(), "bytespool put after by gc") {
		t.Fatalf("log not printed: %s", buf.String())
	}
}
