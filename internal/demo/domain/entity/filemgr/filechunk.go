package filemgr

import (
	"bytes"
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

func msgToBytesWithExtraLen(_ context.Context, magic string, msg msgType, body []byte, extralen int) []byte {
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(magic))
	data[3] = msg.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)+extralen))

	copy(data[8:], body)
	return data
}

func msgToBytes(_ context.Context, magic string, msg msgType, body []byte) []byte {
	data := make([]byte, headerLen+len(body))
	copy(data[0:3], []byte(magic))
	data[3] = msg.Byte()
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(body)))

	copy(data[8:], body)
	return data
}

func fileEventMsgToBytes(ctx context.Context, body []byte) []byte {
	return msgToBytes(ctx, reqMagic, msgFileEvent, body)
}

func fileEventAckToBytes(ctx context.Context, body []byte) []byte {
	return msgToBytes(ctx, ackMagic, msgFileEvent, body)
}

func fileInfoMsgToBytes(ctx context.Context, body []byte) []byte {
	return msgToBytes(ctx, reqMagic, msgFileInfo, body)
}

func fileInfoAckToBytes(ctx context.Context, body []byte) []byte {
	return msgToBytes(ctx, ackMagic, msgFileInfo, body)
}

func fileChunkMsgToBytes(ctx context.Context, body []byte, extralen int) []byte {
	return msgToBytesWithExtraLen(ctx, reqMagic, msgFileChunk, body, extralen)
}

func fileChunkAckToBytes(ctx context.Context, body []byte) []byte {
	return msgToBytes(ctx, ackMagic, msgFileChunk, body)
}

// recvSendFile 处理收到的文件信息，不管是否处理成功，都需要回复给发送方
func recvSendFile(ctx context.Context, body []byte, repo Repository) []byte {
	var sf SendFile
	if err := json.Unmarshal(body, &sf); err != nil {
		g.Log().Errorf(ctx, "recv sendfile fail:%v", err)
		return fileInfoAckToBytes(ctx, []byte("unmarshal fail"))
	}
	var path string
	oldpath := sf.FilePath
	// 检查是否重名，如果重名，那就在文件名后面追加(x)重命名
	path = utils.NextFileName(sf.FilePath, recvedPath)
	g.Log().Debugf(ctx, "file save path: %v", path)

	err := repo.SaveRecvFile(ctx, &RecvFile{
		TaskID:         sf.TaskID,
		TaskName:       sf.TaskName,
		FilePathSave:   path,
		FilePathOrigin: oldpath,
		FileId:         sf.FileId,
		FileSize:       sf.FileSize,
		ChunkNumTotal:  sf.ChunkNumTotal,
		ChunkNumRecved: 0,
		Status:         sf.Status,
	})
	if err != nil {
		g.Log().Errorf(ctx, "save recvfile fail:%v", err)
		return fileInfoAckToBytes(ctx, []byte("save error"))
	}

	return fileInfoAckToBytes(ctx, []byte(gfile.Basename(oldpath)))
}

func recvSendFileChunk(ctx context.Context, body []byte, repo Repository) []byte {
	var sc SendChunk
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&sc); err != nil {
		g.Log().Errorf(ctx, "recv sendChunk fail:%v", err)
		return fileChunkAckToBytes(ctx, []byte("unmarshal fail"))
	}
	// 除去sendchunk结构占用的字节，后面就是文件块数据
	chunkBytes := body[len(body)-sc.ChunkSize:]

	// 为每个收到的文件块创建一个fileSaver实例
	// 从缓存中取出文件块持久化管理服务， 如果是首次存储，将先实例化一个文件接收器：NewFileSave
	fsaver, err := getFileSaver(ctx, sc.FileID, repo)
	if err != nil || fsaver == nil {
		g.Log().Errorf(ctx, "get filesaver fail:%v", err)
		return fileChunkAckToBytes(ctx, []byte("save error"))
	}

	if err := fsaver.SaveChunk(ctx, &fileChunk{
		offset:     sc.ChunkOffset,
		data:       chunkBytes,
		chunkIndex: uint32(sc.ChunkIndex),
		md5:        "",
	}); err != nil {
		g.Log().Errorf(ctx, "save recvfile fail:%v", err)
		return fileChunkAckToBytes(ctx, []byte("save data error"))
	}
	return fileChunkAckToBytes(ctx, []byte(sc.FileID))
}

func recvEvent(ctx context.Context, body []byte, repo Repository) []byte {
	var sc EventMsg
	if err := json.Unmarshal(body, &sc); err != nil {
		g.Log().Errorf(ctx, "recv event fail:%v", err)
		return fileEventAckToBytes(ctx, []byte("unmarshal fail"))
	}
	files, err := repo.GetRecvTask(ctx, sc.TaskId)
	if err != nil {
		g.Log().Errorf(ctx, "recv event fail:%v", err)
		return fileEventAckToBytes(ctx, []byte("recv task fail"))
	}
	for _, file := range files {
		fsaver, err := mustGetFileSaver(ctx, file.FileId)
		if err != nil {
			g.Log().Errorf(ctx, "save recvfile fail:%v", err)
			return fileEventAckToBytes(ctx, []byte("get saver error"))
		}
		if fsaver != nil {
			fsaver.EventNotify(sc.Status)
		}
	}
	return fileEventAckToBytes(ctx, []byte(sc.TaskId))
}
