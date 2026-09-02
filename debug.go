//go:build debug
// +build debug

package bytespool

import (
	"context"
	"log/slog"
	"runtime"
	"strconv"
	"sync"
)

var (
	mu      sync.RWMutex
	records = map[uintptr]*DebugRecord{} // ptr --> *Record
)

func DebugLength() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(records)
}

func DebugClear() {
	mu.Lock()
	clear(records)
	mu.Unlock()
}

func DebugRange(fn func(record *DebugRecord) (next bool)) {
	mu.RLock()
	defer mu.RUnlock()
	for _, r := range records {
		if !fn(r) {
			break
		}
	}
}

type DebugRecord struct {
	Ptr uintptr
	Idx poolIdx
	Pcs [64]uintptr
}

func newRecord(ptr uintptr, idx poolIdx) *DebugRecord {
	r := &DebugRecord{Ptr: ptr, Idx: idx}
	runtime.Callers(2, r.Pcs[:])
	return r
}

func (r *DebugRecord) LogValue() slog.Value {
	fs := runtime.CallersFrames(r.Pcs[:])
	frames := make([]string, 0, 8)
	for {
		f, more := fs.Next()
		frames = append(frames, f.File+":"+strconv.Itoa(f.Line))
		if !more {
			break
		}
	}

	if r.Ptr == 0 && r.Idx == 0 {
		return slog.AnyValue(frames)
	}
	return slog.GroupValue(
		slog.Uint64("ptr", uint64(r.Ptr)),
		slog.Int("size", r.Idx.bytes()),
		slog.Any("stack", slog.AnyValue(frames)),
	)
}

func debug_get(ptr uintptr, idx poolIdx) {
	rec := newRecord(ptr, idx)

	mu.Lock()
	old, loaded := records[ptr]
	records[ptr] = rec
	mu.Unlock()

	if loaded {
		// block was GC'd before being returned to the pool
		log("bytespool put after by gc", slog.Any("record", old))
	}
}

func debug_put(ptr uintptr, idx poolIdx) {
	mu.Lock()
	old, loaded := records[ptr]
	if loaded {
		delete(records, ptr)
	}
	mu.Unlock()

	if loaded {
		if old.Idx != idx {
			// poolIdx in header was corrupted
			log("bytespool put damaged", slog.Any("old", old.Pcs), slog.Any("now", newRecord(0, 0)))
		}
	} else {
		// returning memory not allocated by this pool, magic number collision
		log("bytespool put untracked", slog.Any("stack", newRecord(0, 0)))
	}
}

func debug_log_get_exceed(bytes int) {
	log("bytespool get exceed", slog.Int("bytes", bytes), slog.Any("stack", newRecord(0, 0)))
}

func debug_log_put_nil() {
	log("bytespool put nil", slog.Any("stack", newRecord(0, 0)))
}

func debug_log_put_invalid(ptr uintptr) {
	log("bytespool put invalid", slog.Uint64("ptr", uint64(ptr)), slog.Any("stack", newRecord(0, 0)))
}

var DebugLog func(msg string, attrs ...slog.Attr) = func(msg string, attrs ...slog.Attr) {
	slog.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

func log(msg string, attrs ...slog.Attr) { DebugLog(msg, attrs...) }
