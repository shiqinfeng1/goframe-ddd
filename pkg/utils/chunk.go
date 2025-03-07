package utils

import "github.com/shiqinfeng1/goframe-ddd/pkg/errors"

// 定义不同的字节大小常量
const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

// SplitFile 根据文件大小计算分块大小和起始偏移位置
func SplitFile(size int64) ([]int64, []int, error) {
	// 检查文件大小是否超过 4GB
	if size > 4*GB {
		return nil, nil, errors.ErrOver4GSize
	}

	var chunkSize int
	// 根据文件大小确定分块大小
	switch {
	case size <= 1*MB:
		chunkSize = 1 * MB
	case size <= 100*MB:
		chunkSize = 1 * MB
	case size <= 400*MB:
		chunkSize = 4 * MB
	case size <= 1*GB:
		chunkSize = 8 * MB
	case size <= 4*GB:
		chunkSize = 10 * MB
	}

	var chunkSizes []int
	var chunkOffsets []int64
	offset := int64(0)

	// 计算分块大小和起始偏移位置
	for offset < size {
		// 如果剩余大小小于分块大小，以剩余大小作为当前分块大小
		if size-offset < int64(chunkSize) {
			chunkSizes = append(chunkSizes, int(size-offset))
		} else {
			chunkSizes = append(chunkSizes, chunkSize)
		}
		chunkOffsets = append(chunkOffsets, offset)
		offset += int64(chunkSize)
	}

	return chunkOffsets, chunkSizes, nil
}
