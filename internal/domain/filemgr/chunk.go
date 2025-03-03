package filemgr

import (
	"context"
	"crypto/md5"
	"encoding/binary"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type fileChunk struct {
	fileId     string
	offset     uint64
	data       []byte
	chunkIndex uint32
	md5        string
}

type fileChunkHeaderReq struct {
	Magic       string // len=4
	FileId      string // len=20
	ChunkIdx    uint32 // len=4
	ChunkOffset uint64 // len=8
	ChunkSize   uint64 // len=8
	Md5         string // len=32
}
type fileChunkHeaderResp struct {
	Magic    string // len=4
	FileId   string // len=20
	ChunkIdx uint32 // len=4
	Status   uint32 // len=4
}

func NewFileChunkHeader(b []byte) *fileChunkHeaderReq {
	return &fileChunkHeaderReq{
		Magic:       gconv.String(b[:4]),
		FileId:      gconv.String(b[4:24]),
		ChunkIdx:    binary.LittleEndian.Uint32(b[24:28]),
		ChunkOffset: binary.LittleEndian.Uint64(b[28:36]),
		ChunkSize:   binary.LittleEndian.Uint64(b[36:44]),
		Md5:         gconv.String(b[44 : 44+md5.Size*2]),
	}
}

func ToFileChunkHeaderReq(fileId, hash string, chunkIdx int, chunkOffset, chunkSize int64) []byte {
	b := make([]byte, headerLenReq)
	copy(b[:4], sendMagic)
	copy(b[4:24], fileId)
	binary.LittleEndian.PutUint32(b[24:28], uint32(chunkIdx))
	binary.LittleEndian.PutUint64(b[28:36], uint64(chunkOffset))
	binary.LittleEndian.PutUint64(b[36:44], uint64(chunkSize))
	copy(b[44:44+md5.Size*2], hash)
	return b
}

func ToFileChunkHeaderResp(fileId string, chunkIdx, status int) []byte {
	b := make([]byte, headerLenReq)
	copy(b[:4], sendMagic)
	copy(b[4:24], fileId)
	binary.LittleEndian.PutUint32(b[24:28], uint32(chunkIdx))
	binary.LittleEndian.PutUint32(b[28:32], uint32(status))
	return b
}

// 计算文件块信息
// 返回：文件块索引列表，每个块的偏移列表，每个文件块的大小
func CalcChunks(ctx context.Context, totalSize int64) ([]int64, []int64) {
	// 块大小默认16M
	chunkSize := g.Cfg().MustGet(ctx, "chunksize", 16).Int64()
	chunkSize = chunkSize * 1024 * 1024
	offsets := make([]int64, 0)
	sizes := make([]int64, 0)

	for offset := int64(0); offset < totalSize; offset += chunkSize {
		offsets = append(offsets, offset)
		if offset+chunkSize > totalSize {
			sizes = append(sizes, totalSize-offset)
		} else {
			sizes = append(sizes, chunkSize)
		}
	}
	return offsets, sizes
}
