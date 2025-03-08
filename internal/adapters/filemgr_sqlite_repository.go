package adapters

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent"
	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
)

type filemgrRepo struct {
	db *ent.Client
}

// NewTrainingRepo .
func NewFilemgrRepo(db *ent.Client) *filemgrRepo {
	return &filemgrRepo{
		db: db,
	}
}

func (*filemgrRepo) GetSendFile(ctx context.Context, taskId, filePath string) (*filemgr.SendFile, error) {
	return nil, nil
}

func (*filemgrRepo) GetSendTask(ctx context.Context, taskId string) ([]*filemgr.SendFile, error) {
	return nil, nil
}

func (*filemgrRepo) SaveSendFile(ctx context.Context, sendFile *filemgr.SendFile) (int, error) {
	return 0, nil
}

// 插入sendchunk和更新sendfile的chunk统计
func (*filemgrRepo) UpdateSendChunk(ctx context.Context, sendChunk *filemgr.SendChunk) error {
	return nil
}

func (*filemgrRepo) GetRecvTask(ctx context.Context, taskId string) ([]*filemgr.RecvFile, error) {
	return nil, nil
}

func (*filemgrRepo) SaveRecvFile(ctx context.Context, recvile *filemgr.RecvFile) error {
	return nil
}

// 插入sendchunk和更新sendfile的chunk统计
func (*filemgrRepo) UpdateRecvChunk(ctx context.Context, recvChunk *filemgr.RecvChunk) error {
	return nil
}

func (*filemgrRepo) GetRecvTaskFile(ctx context.Context, taskId, fileId string) (*filemgr.RecvFile, error) {
	return nil, nil
}

func (*filemgrRepo) CountOfRecvedChunks(ctx context.Context, taskId, fileId string) (int, error) {
	return 0, nil
}
