package adapters

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/recvfile"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/sendfile"
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

func (f *filemgrRepo) GetSendFile(ctx context.Context, taskId, filePath string) (*filemgr.SendFile, error) {
	sf, err := f.db.SendFile.
		Query().
		Where(sendfile.TaskID(taskId), sendfile.FilePath(filePath)).
		Only(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query sandfile fail")
	}
	out := &filemgr.SendFile{
		ID:             sf.ID,
		TaskID:         sf.TaskID,
		TaskName:       sf.TaskName,
		FilePath:       sf.FilePath,
		FileId:         sf.FileID,
		FileSize:       sf.FileSize,
		ChunkNumTotal:  sf.ChunkNumTotal,
		ChunkNumSended: sf.ChunkNumSended,
		Status:         sf.Status,
		Elapsed:        sf.Elapsed,
		Speed:          sf.Speed,
	}
	return out, nil
}

func (f *filemgrRepo) GetSendTask(ctx context.Context, taskId string) ([]*filemgr.SendFile, error) {
	sfs, err := f.db.SendFile.
		Query().
		Where(sendfile.TaskID(taskId)).
		All(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query sandfile fail")
	}
	var out []*filemgr.SendFile
	for _, sf := range sfs {
		out = append(out, &filemgr.SendFile{
			ID:             sf.ID,
			TaskID:         sf.TaskID,
			TaskName:       sf.TaskName,
			FilePath:       sf.FilePath,
			FileId:         sf.FileID,
			FileSize:       sf.FileSize,
			ChunkNumTotal:  sf.ChunkNumTotal,
			ChunkNumSended: sf.ChunkNumSended,
			Status:         sf.Status,
			Elapsed:        sf.Elapsed,
			Speed:          sf.Speed,
		})
	}
	return out, nil
}

func (f *filemgrRepo) SaveSendFile(ctx context.Context, sf *filemgr.SendFile) (int, error) {
	created, err := f.db.SendFile.
		Create().
		SetTaskID(sf.TaskID).
		SetTaskName(sf.TaskName).
		SetFilePath(sf.FilePath).
		SetFileID(sf.FileId).
		SetFileSize(sf.FileSize).
		SetChunkNumTotal(sf.ChunkNumTotal).
		SetChunkNumSended(sf.ChunkNumSended).
		SetStatus(sf.Status).
		SetElapsed(sf.Elapsed).
		SetSpeed(sf.Speed).
		Save(ctx)
	if err != nil {
		return 0, gerror.Wrap(err, "save sendfile fail")
	}
	return created.ID, nil
}

// 插入sendchunk和更新sendfile的chunk统计
func (f *filemgrRepo) UpdateSendChunk(ctx context.Context, sc *filemgr.SendChunk) error {
	// 开始事务
	tx, err := f.db.Tx(ctx)
	if err != nil {
		return err
	}
	sf, err := tx.SendFile.
		Query().
		Where(sendfile.FileID(sc.FileID)).
		Only(ctx)
	if err != nil {
		return tx.Rollback()
	}
	// 插入filechunk记录
	_, err = tx.SendChunk.
		Create().
		SetChunkIndex(sc.ChunkIndex).
		SetChunkOffset(sc.ChunkOffset).
		SetChunkSize(sc.ChunkSize).
		SetSendFileID(sf.ID).
		Save(ctx)
	if err != nil {
		return tx.Rollback()
	}
	var status int
	if sf.ChunkNumTotal == sf.ChunkNumSended+1 {
		status = filemgr.StatusSuccessful.Int()
	} else {
		status = filemgr.StatusSending.Int()
	}
	_, err = tx.SendFile.
		UpdateOneID(sf.ID).
		AddChunkNumSended(1).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return tx.Rollback()
	}

	// 提交事务C
	return tx.Commit()
}

func (f *filemgrRepo) GetRecvTask(ctx context.Context, taskId string) ([]*filemgr.RecvFile, error) {
	sfs, err := f.db.RecvFile.
		Query().
		Where(recvfile.TaskID(taskId)).
		All(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query sandfile fail")
	}
	var out []*filemgr.RecvFile
	for _, sf := range sfs {
		out = append(out, &filemgr.RecvFile{
			TaskID:         sf.TaskID,
			TaskName:       sf.TaskName,
			FilePathSave:   sf.FilePathSave,
			FilePathOrigin: sf.FilePathOrigin,
			FileId:         sf.FileID,
			FileSize:       sf.FileSize,
			ChunkNumTotal:  sf.ChunkNumTotal,
			ChunkNumRecved: sf.ChunkNumRecved,
			Status:         sf.Status,
		})
	}
	return out, nil
}

func (f *filemgrRepo) SaveRecvFile(ctx context.Context, rf *filemgr.RecvFile) error {
	_, err := f.db.RecvFile.
		Create().
		SetTaskID(rf.TaskID).
		SetTaskName(rf.TaskName).
		SetFilePathOrigin(rf.FilePathOrigin).
		SetFilePathSave(rf.FilePathSave).
		SetFileID(rf.FileId).
		SetFileSize(rf.FileSize).
		SetChunkNumTotal(rf.ChunkNumTotal).
		SetChunkNumRecved(rf.ChunkNumRecved).
		SetStatus(rf.Status).
		Save(ctx)
	if err != nil {
		return gerror.Wrap(err, "save sendfile fail")
	}
	return nil
}

// 插入sendchunk和更新sendfile的chunk统计
func (f *filemgrRepo) UpdateRecvChunk(ctx context.Context, rc *filemgr.RecvChunk) (bool, error) {
	// 开始事务
	tx, err := f.db.Tx(ctx)
	if err != nil {
		return false, err
	}
	rf, err := tx.RecvFile.
		Query().
		Where(recvfile.FileID(rc.FileID)).
		Only(ctx)
	if err != nil {
		return false, tx.Rollback()
	}
	// 插入filechunk记录
	_, err = tx.RecvChunk.
		Create().
		SetChunkIndex(rc.ChunkIndex).
		SetChunkOffset(rc.ChunkOffset).
		SetChunkSize(rc.ChunkSize).
		SetRecvFileID(rf.ID).
		Save(ctx)
	if err != nil {
		return false, tx.Rollback()
	}
	var (
		status   int
		finished bool
	)
	if rf.ChunkNumTotal == rf.ChunkNumRecved+1 {
		status = filemgr.StatusSuccessful.Int()
		finished = true
	} else {
		status = filemgr.StatusSending.Int()
	}
	_, err = tx.RecvFile.
		UpdateOneID(rf.ID).
		AddChunkNumRecved(1).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return false, tx.Rollback()
	}

	// 提交事务C
	return finished, tx.Commit()
}

func (f *filemgrRepo) GetRecvFile(ctx context.Context, fileId string) (*filemgr.RecvFile, error) {
	rf, err := f.db.RecvFile.
		Query().
		Where(recvfile.FileID(fileId)).
		Only(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query sandfile fail")
	}
	out := &filemgr.RecvFile{
		TaskID:         rf.TaskID,
		TaskName:       rf.TaskName,
		FilePathSave:   rf.FilePathSave,
		FilePathOrigin: rf.FilePathOrigin,
		FileId:         rf.FileID,
		FileSize:       rf.FileSize,
		ChunkNumTotal:  rf.ChunkNumTotal,
		ChunkNumRecved: rf.ChunkNumRecved,
		Status:         rf.Status,
	}
	return out, nil
}

func (f *filemgrRepo) CountOfRecvedChunks(ctx context.Context, fileId string) (int, error) {
	return 0, nil
}

func (f *filemgrRepo) UpdateSendStatus(ctx context.Context, fileId string, status filemgr.Status) error {
	_, err := f.db.SendFile.
		Update().
		Where(sendfile.FileID(fileId)).
		SetStatus(status.Int()).
		Save(ctx)
	if err != nil {
		return gerror.Wrap(err, "update sendfile status fail")
	}
	return nil
}
