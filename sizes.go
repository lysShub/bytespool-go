package bytespool

import (
	"sync"
	"unsafe"
)

const (
	maxPoolIdx    = poolIdx(len(pools)) - 1
	poolIdxExceed = maxPoolIdx + 1

	smallSizeMax  = 1024
	middleSizeMax = 32768
	smallSizeDiv  = 8
	middleSizeDiv = 128

	large64k   = 64 << 10
	large128k  = 128 << 10
	large1024k = 1 << 20
	large4096k = 4 << 20
)

var pools = [70]*sync.Pool{
	0:  {New: func() any { return unsafe.Pointer(new([16]byte)) }},
	1:  {New: func() any { return unsafe.Pointer(new([24]byte)) }},
	2:  {New: func() any { return unsafe.Pointer(new([32]byte)) }},
	3:  {New: func() any { return unsafe.Pointer(new([48]byte)) }},
	4:  {New: func() any { return unsafe.Pointer(new([64]byte)) }},
	5:  {New: func() any { return unsafe.Pointer(new([80]byte)) }},
	6:  {New: func() any { return unsafe.Pointer(new([96]byte)) }},
	7:  {New: func() any { return unsafe.Pointer(new([112]byte)) }},
	8:  {New: func() any { return unsafe.Pointer(new([128]byte)) }},
	9:  {New: func() any { return unsafe.Pointer(new([144]byte)) }},
	10: {New: func() any { return unsafe.Pointer(new([160]byte)) }},
	11: {New: func() any { return unsafe.Pointer(new([176]byte)) }},
	12: {New: func() any { return unsafe.Pointer(new([192]byte)) }},
	13: {New: func() any { return unsafe.Pointer(new([208]byte)) }},
	14: {New: func() any { return unsafe.Pointer(new([224]byte)) }},
	15: {New: func() any { return unsafe.Pointer(new([240]byte)) }},
	16: {New: func() any { return unsafe.Pointer(new([256]byte)) }},
	17: {New: func() any { return unsafe.Pointer(new([288]byte)) }},
	18: {New: func() any { return unsafe.Pointer(new([320]byte)) }},
	19: {New: func() any { return unsafe.Pointer(new([352]byte)) }},
	20: {New: func() any { return unsafe.Pointer(new([384]byte)) }},
	21: {New: func() any { return unsafe.Pointer(new([416]byte)) }},
	22: {New: func() any { return unsafe.Pointer(new([448]byte)) }},
	23: {New: func() any { return unsafe.Pointer(new([480]byte)) }},
	24: {New: func() any { return unsafe.Pointer(new([512]byte)) }},
	25: {New: func() any { return unsafe.Pointer(new([576]byte)) }},
	26: {New: func() any { return unsafe.Pointer(new([640]byte)) }},
	27: {New: func() any { return unsafe.Pointer(new([704]byte)) }},
	28: {New: func() any { return unsafe.Pointer(new([768]byte)) }},
	29: {New: func() any { return unsafe.Pointer(new([896]byte)) }},
	30: {New: func() any { return unsafe.Pointer(new([1024]byte)) }},
	31: {New: func() any { return unsafe.Pointer(new([1152]byte)) }},
	32: {New: func() any { return unsafe.Pointer(new([1280]byte)) }},
	33: {New: func() any { return unsafe.Pointer(new([1408]byte)) }},
	34: {New: func() any { return unsafe.Pointer(new([1536]byte)) }},
	35: {New: func() any { return unsafe.Pointer(new([1792]byte)) }},
	36: {New: func() any { return unsafe.Pointer(new([2048]byte)) }},
	37: {New: func() any { return unsafe.Pointer(new([2304]byte)) }},
	38: {New: func() any { return unsafe.Pointer(new([2688]byte)) }},
	39: {New: func() any { return unsafe.Pointer(new([3072]byte)) }},
	40: {New: func() any { return unsafe.Pointer(new([3200]byte)) }},
	41: {New: func() any { return unsafe.Pointer(new([3456]byte)) }},
	42: {New: func() any { return unsafe.Pointer(new([4096]byte)) }},
	43: {New: func() any { return unsafe.Pointer(new([4864]byte)) }},
	44: {New: func() any { return unsafe.Pointer(new([5376]byte)) }},
	45: {New: func() any { return unsafe.Pointer(new([6144]byte)) }},
	46: {New: func() any { return unsafe.Pointer(new([6528]byte)) }},
	47: {New: func() any { return unsafe.Pointer(new([6784]byte)) }},
	48: {New: func() any { return unsafe.Pointer(new([6912]byte)) }},
	49: {New: func() any { return unsafe.Pointer(new([8192]byte)) }},
	50: {New: func() any { return unsafe.Pointer(new([9472]byte)) }},
	51: {New: func() any { return unsafe.Pointer(new([9728]byte)) }},
	52: {New: func() any { return unsafe.Pointer(new([10240]byte)) }},
	53: {New: func() any { return unsafe.Pointer(new([10880]byte)) }},
	54: {New: func() any { return unsafe.Pointer(new([12288]byte)) }},
	55: {New: func() any { return unsafe.Pointer(new([13568]byte)) }},
	56: {New: func() any { return unsafe.Pointer(new([14336]byte)) }},
	57: {New: func() any { return unsafe.Pointer(new([16384]byte)) }},
	58: {New: func() any { return unsafe.Pointer(new([18432]byte)) }},
	59: {New: func() any { return unsafe.Pointer(new([19072]byte)) }},
	60: {New: func() any { return unsafe.Pointer(new([20480]byte)) }},
	61: {New: func() any { return unsafe.Pointer(new([21760]byte)) }},
	62: {New: func() any { return unsafe.Pointer(new([24576]byte)) }},
	63: {New: func() any { return unsafe.Pointer(new([27264]byte)) }},
	64: {New: func() any { return unsafe.Pointer(new([28672]byte)) }},
	65: {New: func() any { return unsafe.Pointer(new([32768]byte)) }},
	66: {New: func() any { return unsafe.Pointer(new([large64k]byte)) }},
	67: {New: func() any { return unsafe.Pointer(new([large128k]byte)) }},
	68: {New: func() any { return unsafe.Pointer(new([large1024k]byte)) }},
	69: {New: func() any { return unsafe.Pointer(new([large4096k]byte)) }},
}

type poolIdx uint8 //
var (
	classToSize    = [len(pools)]uint32{16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512, 576, 640, 704, 768, 896, 1024, 1152, 1280, 1408, 1536, 1792, 2048, 2304, 2688, 3072, 3200, 3456, 4096, 4864, 5376, 6144, 6528, 6784, 6912, 8192, 9472, 9728, 10240, 10880, 12288, 13568, 14336, 16384, 18432, 19072, 20480, 21760, 24576, 27264, 28672, 32768, 64 << 10, 128 << 10, 1 << 20, 4 << 20}
	sizeToClass8   = [smallSizeMax/smallSizeDiv + 1]uint8{0, 1, 2, 3, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17, 17, 18, 18, 19, 19, 19, 19, 20, 20, 20, 20, 21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24, 24, 24, 24, 25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 27, 27, 27, 27, 28, 28, 28, 28, 28, 28, 28, 28, 29, 29, 29, 29, 29, 29, 29, 29, 30, 30, 30, 30, 30, 30, 30, 30, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
	sizeToClass128 = [(middleSizeMax-smallSizeMax)/middleSizeDiv + 1]uint8{32, 33, 34, 35, 36, 37, 37, 38, 38, 39, 39, 40, 40, 40, 41, 41, 41, 42, 43, 43, 44, 44, 44, 44, 44, 45, 45, 45, 45, 45, 45, 46, 46, 46, 46, 47, 47, 47, 47, 47, 47, 48, 48, 48, 49, 49, 50, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 52, 52, 52, 52, 52, 52, 52, 52, 52, 52, 53, 53, 54, 54, 54, 54, 55, 55, 55, 55, 55, 56, 56, 56, 56, 56, 56, 56, 56, 56, 56, 56, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 58, 58, 58, 58, 58, 58, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 61, 61, 61, 61, 61, 62, 62, 62, 62, 62, 62, 62, 62, 62, 62, 62, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67}
)

func newPoolIdx(bytes int) poolIdx {
	if Debug && bytes < hdrsize {
		panic(bytes)
	}
	if bytes <= middleSizeMax {
		var i int
		if bytes <= smallSizeMax {
			i = int(sizeToClass8[(bytes+smallSizeDiv-1)>>3])
		} else {
			i = int(sizeToClass128[(bytes-smallSizeMax+middleSizeDiv-1)>>7])
		}
		i = max(i, 2)
		return poolIdx(i - 2)
	}
	switch {
	case bytes <= large64k:
		return 66
	case bytes <= large128k:
		return 67
	case bytes <= large1024k:
		return 68
	case bytes <= large4096k:
		return 69
	}
	return poolIdxExceed
}
func (p poolIdx) bytes() int {
	if Debug && p >= poolIdxExceed {
		panic(p)
	}
	return int(classToSize[p])
}
