//go:build !debug
// +build !debug

package bytespool

const Debug = false

type Record struct{}

func DebugLength() int                               { return 0 }
func DebugClear()                                    {}
func DebugRange(fn func(record *Record) (next bool)) {}

func debug_get(ptr uintptr, idx poolIdx) {}
func debug_put(ptr uintptr, idx poolIdx) {}
func debug_log_get_exceed(bytes int)     {}
func debug_log_put_nil()                 {}
func debug_log_put_invalid(ptr uintptr)  {}
