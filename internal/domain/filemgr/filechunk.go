package filemgr

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
)

var recvedPath = ""

func init() {
	home, _ := gfile.Home()
	recvedPath = filepath.Join(home, "Downloads")
}

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

func fileInfoMsgToBytes(_ context.Context, body []byte) []byte {
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(reqMagic))
	data[3] = msgFileInfo.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)))

	copy(data[8:], body)
	return data
}

func fileChunkMsgToBytes(_ context.Context, body []byte) []byte {
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(reqMagic))
	data[3] = msgFileChunk.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)))

	copy(data[8:], body)
	return data
}

func fileInfoAckToBytes(_ context.Context, body []byte) []byte {
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(ackMagic))
	data[3] = msgFileInfo.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)))

	copy(data[8:], body)
	return data
}

func fileChunkAckToBytes(_ context.Context, body []byte) []byte {
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(ackMagic))
	data[3] = msgFileChunk.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)))

	copy(data[8:], body)
	return data
}

// recvSendFile 处理收到的文件信息，不管是否处理成功，都需要回复给发送方
func recvSendFile(ctx context.Context, body []byte, repo Repository) []byte {
	var sendFile SendFile
	if err := json.Unmarshal(body, &sendFile); err != nil {
		g.Log().Errorf(ctx, "recv sendfile fail:%v", err)
		return fileInfoAckToBytes(ctx, []byte("unmarshal fail"))
	}
	var path string
	oldpath := sendFile.FilePath
	// 检查是否重名，如果重名，那就在文件名后面追加(x)重命名
	path = utils.NextFileName(sendFile.FilePath, recvedPath)
	g.Log().Debugf(ctx, "file save path: %v", path)

	err := repo.SaveRecvFile(ctx, &RecvFile{
		TaskID:         sendFile.TaskID,
		TaskName:       sendFile.TaskName,
		FilePathSave:   path,
		FilePathOrigin: oldpath,
		Fid:            sendFile.Fid,
		FileSize:       sendFile.FileSize,
		ChunkNumTotal:  sendFile.ChunkNumTotal,
		ChunkNumRecved: 0,
		Status:         0,
	})
	if err != nil {
		g.Log().Errorf(ctx, "save recvfile fail:%v", err)
		return fileInfoAckToBytes(ctx, []byte("save error"))
	}

	return fileInfoAckToBytes(ctx, []byte(path))
}

func recvSendFileChunk(ctx context.Context, body []byte, repo Repository) []byte {
	var sendChunk SendChunk
	if err := json.Unmarshal(body, &sendChunk); err != nil {
		g.Log().Errorf(ctx, "recv sendChunk fail:%v", err)
		return fileChunkAckToBytes(ctx, []byte("unmarshal fail"))
	}
	// 除去sendchunk结构占用的字节，后面就是文件块数据
	chunkBytes := body[len(body)-sendChunk.ChunkSize:]

	err := repo.UpdateRecvChunk(ctx, &RecvChunk{
		FileID:      0,
		ChunkIndex:  sendChunk.ChunkIndex,
		ChunkOffset: sendChunk.ChunkOffset,
		ChunkSize:   sendChunk.ChunkSize,
		Status:      0,
	})
	if err != nil {
		g.Log().Errorf(ctx, "save recvfile fail:%v", err)
		return fileChunkAckToBytes(ctx, []byte("save error"))
	}

	return fileChunkAckToBytes(ctx, []byte(path))
}
