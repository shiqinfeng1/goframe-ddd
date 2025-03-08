package filemgr

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/shiqinfeng1/goframe-ddd/pkg/cache"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
)

var recvedPath = ""

func init() {
	home, _ := gfile.Home()
	recvedPath = filepath.Join(home, "Downloads")
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
		FileId:         sendFile.FileId,
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
	var sc SendChunk
	if err := json.Unmarshal(body, &sc); err != nil {
		g.Log().Errorf(ctx, "recv sendChunk fail:%v", err)
		return fileChunkAckToBytes(ctx, []byte("unmarshal fail"))
	}
	// 除去sendchunk结构占用的字节，后面就是文件块数据
	chunkBytes := body[len(body)-sc.ChunkSize:]

	// 为每个收到的文件块创建一个fileSaver实例
	// 从缓存中取出文件块持久化管理服务， 如果是首次存储，将先实例化一个文件接收器：NewFileSave
	fs, err := NewFileSave(ctx, cache.Memory(), sc.TaskID, sc.FileID, repo)
	if err != nil {
		g.Log().Errorf(ctx, "save recvfile fail:%v", err)
		return fileChunkAckToBytes(ctx, []byte("save error"))
	}

	fs.SaveChunk(&fileChunk{
		taskId:     sc.TaskID,
		fileId:     sc.FileID,
		offset:     sc.ChunkOffset,
		data:       chunkBytes,
		chunkIndex: uint32(sc.ChunkIndex),
		md5:        "",
	})

	if err := repo.UpdateRecvChunk(ctx, &RecvChunk{
		TaskID:      sc.TaskID,
		FileID:      sc.FileID,
		ChunkIndex:  sc.ChunkIndex,
		ChunkOffset: sc.ChunkOffset,
		ChunkSize:   sc.ChunkSize,
	}); err != nil {
		g.Log().Errorf(ctx, "save recvfile fail:%v", err)
		return fileChunkAckToBytes(ctx, []byte("save error"))
	}

	return fileChunkAckToBytes(ctx, []byte(sc.FileID))
}
