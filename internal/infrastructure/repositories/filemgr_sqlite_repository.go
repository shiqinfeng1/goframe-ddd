package repositories

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/internal/domain/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories/ent"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories/ent/filetransfertask"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories/ent/recvchunk"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories/ent/recvfile"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories/ent/sendchunk"
	"github.com/shiqinfeng1/goframe-ddd/internal/infrastructure/repositories/ent/sendfile"
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
	if exist, err := f.db.FileTransferTask.
		Query().
		Where(filetransfertask.TaskIDEQ(taskId)).
		Exist(ctx); !exist {
		return nil, gerror.Wrap(err, "query task fail")
	}
	task, err := f.db.FileTransferTask.
		Query().
		Where(filetransfertask.TaskIDEQ(taskId)).
		Only(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query sendfile fail")
	}

	if exist, _ := f.db.SendFile.
		Query().
		Where(sendfile.TaskIDEQ(taskId), sendfile.FilePathEQ(filePath)).
		Exist(ctx); !exist {
		return nil, nil
	}
	sf, err := f.db.SendFile.
		Query().
		Where(sendfile.TaskIDEQ(taskId), sendfile.FilePathEQ(filePath)).
		Only(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query sendfile fail")
	}
	out := &filemgr.SendFile{
		ID:             sf.ID,
		TaskID:         sf.TaskID,
		TaskName:       task.TaskName,
		FilePath:       sf.FilePath,
		FileId:         sf.FileID,
		FileSize:       sf.FileSize,
		ChunkNumTotal:  sf.ChunkNumTotal,
		ChunkNumSended: sf.ChunkNumSended,
		Status:         sf.Status,
	}
	return out, nil
}

func (f *filemgrRepo) RemoveTasks(ctx context.Context, taskids []string) error {
	// 开始事务
	tx, err := f.db.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.FileTransferTask.Delete().Where(filetransfertask.TaskIDIn(taskids...)).Exec(ctx)
	if err != nil {
		return tx.Rollback()
	}
	ids, err := tx.SendFile.Query().Where(sendfile.TaskIDIn(taskids...)).IDs(ctx)
	if err != nil {
		return tx.Rollback()
	}
	_, err = tx.SendChunk.Delete().Where(sendchunk.SendfileIDIn(ids...)).Exec(ctx)
	if err != nil {
		return tx.Rollback()
	}
	_, err = tx.SendFile.Delete().Where(sendfile.TaskIDIn(taskids...)).Exec(ctx)
	if err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}

func (f *filemgrRepo) GetNotCompletedTasks(ctx context.Context) ([]*filemgr.FileTransferTask, map[string][]*filemgr.SendFile, error) {
	exist, _ := f.db.FileTransferTask.
		Query().
		Where(filetransfertask.StatusIn(
			filemgr.StatusWaiting.Int(),
			filemgr.StatusSending.Int(),
			filemgr.StatusPaused.Int(),
			filemgr.StatusFailed.Int())).Exist(ctx)
	if !exist {
		return []*filemgr.FileTransferTask{}, map[string][]*filemgr.SendFile{}, nil
	}

	tasks, err := f.db.FileTransferTask.
		Query().
		Where(filetransfertask.StatusIn(
			filemgr.StatusWaiting.Int(),
			filemgr.StatusSending.Int(),
			filemgr.StatusPaused.Int(),
			filemgr.StatusFailed.Int())).All(ctx)
	if err != nil {
		return nil, nil, gerror.Wrap(err, "query filetransfertask fail")
	}
	out := make(map[string][]*filemgr.SendFile)
	outTasks := make([]*filemgr.FileTransferTask, 0)

	for _, task := range tasks {
		outTasks = append(outTasks, &filemgr.FileTransferTask{
			TaskID:   task.TaskID,
			TaskName: task.TaskName,
			NodeID:   task.NodeID,
			Elapsed:  task.Elapsed,
			Speed:    task.Speed,
			Status:   task.Status,
		})

		out[task.TaskID] = make([]*filemgr.SendFile, 0)
		sfs, err := f.GetSendFilesByTask(ctx, task.TaskID)
		if err != nil {
			return nil, nil, err
		}
		out[task.TaskID] = append(out[task.TaskID], sfs...)
	}

	return outTasks, out, nil
}

func (f *filemgrRepo) GetCompletedTasks(ctx context.Context) ([]*filemgr.FileTransferTask, map[string][]*filemgr.SendFile, error) {
	exist, _ := f.db.FileTransferTask.
		Query().
		Where(filetransfertask.StatusEQ(filemgr.StatusSuccessful.Int())).Exist(ctx)
	if !exist {
		return []*filemgr.FileTransferTask{}, map[string][]*filemgr.SendFile{}, nil
	}

	tasks, err := f.db.FileTransferTask.
		Query().
		Where(filetransfertask.StatusEQ(filemgr.StatusSuccessful.Int())).All(ctx)
	if err != nil {
		return nil, nil, gerror.Wrap(err, "query filetransfertask fail")
	}
	out := make(map[string][]*filemgr.SendFile)
	outTasks := make([]*filemgr.FileTransferTask, 0)

	for _, task := range tasks {
		outTasks = append(outTasks, &filemgr.FileTransferTask{
			TaskID:   task.TaskID,
			TaskName: task.TaskName,
			NodeID:   task.NodeID,
			Elapsed:  task.Elapsed,
			Speed:    task.Speed,
			Status:   task.Status,
		})

		out[task.TaskID] = make([]*filemgr.SendFile, 0)
		sfs, err := f.GetSendFilesByTask(ctx, task.TaskID)
		if err != nil {
			return nil, nil, err
		}
		out[task.TaskID] = append(out[task.TaskID], sfs...)
	}

	return outTasks, out, nil
}

func (f *filemgrRepo) GetSendFilesByTask(ctx context.Context, taskId string) ([]*filemgr.SendFile, error) {
	exist, _ := f.db.SendFile.
		Query().
		Where(sendfile.TaskIDEQ(taskId)).Exist(ctx)
	if !exist {
		return []*filemgr.SendFile{}, nil
	}
	sfs, err := f.db.SendFile.
		Query().
		Where(sendfile.TaskIDEQ(taskId)).All(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query sendfile fail")
	}

	var out []*filemgr.SendFile
	for _, sf := range sfs {
		out = append(out, &filemgr.SendFile{
			ID:             sf.ID,
			TaskID:         sf.TaskID,
			FilePath:       sf.FilePath,
			FileId:         sf.FileID,
			FileSize:       sf.FileSize,
			ChunkNumTotal:  sf.ChunkNumTotal,
			ChunkNumSended: sf.ChunkNumSended,
			Status:         sf.Status,
		})
	}
	return out, nil
}

func (f *filemgrRepo) SaveTask(ctx context.Context, ftt *filemgr.FileTransferTask) error {
	_, err := f.db.FileTransferTask.
		Create().
		SetTaskID(ftt.TaskID).
		SetTaskName(ftt.TaskName).
		SetNodeID(ftt.NodeID).
		SetStatus(ftt.Status).
		Save(ctx)
	if err != nil {
		return gerror.Wrap(err, "save task fail")
	}
	return nil
}

func (f *filemgrRepo) SaveSendFile(ctx context.Context, sf *filemgr.SendFile) (int, error) {
	created, err := f.db.SendFile.
		Create().
		SetTaskID(sf.TaskID).
		SetFilePath(sf.FilePath).
		SetFileID(sf.FileId).
		SetFileSize(sf.FileSize).
		SetChunkNumTotal(sf.ChunkNumTotal).
		SetChunkNumSended(sf.ChunkNumSended).
		SetStatus(sf.Status).
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
		Where(sendfile.FileIDEQ(sc.FileID)).
		Only(ctx)
	if err != nil {
		tx.Rollback()
		return gerror.Wrapf(err, "query sendfile fail")
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
		tx.Rollback()
		return gerror.Wrapf(err, "save chunk fail")
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
		tx.Rollback()
		return gerror.Wrapf(err, "update sendfile fail")
	}
	_, err = tx.FileTransferTask.
		Update().
		Where(filetransfertask.TaskIDEQ(sf.TaskID)).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return gerror.Wrapf(err, "update task status fail")
	}
	// 提交事务C
	return tx.Commit()
}

func (f *filemgrRepo) GetRecvTask(ctx context.Context, taskId string) ([]*filemgr.RecvFile, error) {
	if exist, _ := f.db.RecvFile.
		Query().
		Where(recvfile.TaskIDEQ(taskId)).Exist(ctx); !exist {
		return []*filemgr.RecvFile{}, nil
	}
	sfs, err := f.db.RecvFile.
		Query().
		Where(recvfile.TaskIDEQ(taskId)).All(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query recvfile fail")
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
		return gerror.Wrap(err, "save recvfile fail")
	}
	return nil
}

// 插入sendchunk和更新sendfile的chunk统计
func (f *filemgrRepo) UpdateRecvChunk(ctx context.Context, rc *filemgr.RecvChunk) (*filemgr.RecvFile, error) {
	// 开始事务
	tx, err := f.db.Tx(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "open tx fail")
	}
	rf, err := tx.RecvFile.
		Query().
		Where(recvfile.FileIDEQ(rc.FileID)).
		Only(ctx)
	if err != nil {
		tx.Rollback()
		return nil, gerror.Wrap(err, "query recvfile fail")
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
		tx.Rollback()
		return nil, gerror.Wrap(err, "query recvchunk fail")
	}
	var status int
	if rf.ChunkNumTotal == rf.ChunkNumRecved+1 {
		status = filemgr.StatusSuccessful.Int()
	} else {
		status = filemgr.StatusSending.Int()
	}
	newrf, err := tx.RecvFile.
		UpdateOneID(rf.ID).
		AddChunkNumRecved(1).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, gerror.Wrap(err, "save recvfile fail")
	}

	out := &filemgr.RecvFile{
		TaskID:         newrf.TaskID,
		TaskName:       newrf.TaskName,
		FilePathSave:   newrf.FilePathSave,
		FilePathOrigin: newrf.FilePathOrigin,
		FileId:         newrf.FileID,
		FileSize:       newrf.FileSize,
		ChunkNumTotal:  newrf.ChunkNumTotal,
		ChunkNumRecved: newrf.ChunkNumRecved,
		Status:         newrf.Status,
	}
	// 提交事务C
	return out, tx.Commit()
}

func (f *filemgrRepo) GetRecvFile(ctx context.Context, fileId string) (*filemgr.RecvFile, error) {
	if exist, _ := f.db.RecvFile.Query().
		Where(recvfile.FileIDEQ(fileId)).Exist(ctx); !exist {
		return &filemgr.RecvFile{}, nil
	}
	rf, err := f.db.RecvFile.
		Query().
		Where(recvfile.FileIDEQ(fileId)).Only(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "query recvfile fail")
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
	id, err := f.db.RecvFile.
		Query().
		Where(recvfile.FileIDEQ(fileId)).OnlyID(ctx)
	if err != nil {
		return 0, gerror.Wrap(err, "query recvfile fail")
	}
	cnt, err := f.db.RecvChunk.
		Query().
		Where(recvchunk.RecvfileIDEQ(id)).
		Count(ctx)
	if err != nil {
		return 0, gerror.Wrap(err, "count recvchunk fail")
	}
	return cnt, nil
}

func (f *filemgrRepo) UpdateSpeed(ctx context.Context, taskid, elapsed, speed string) error {
	_, err := f.db.FileTransferTask.
		Update().
		Where(filetransfertask.TaskIDEQ(taskid)).
		SetElapsed(elapsed).
		SetSpeed(speed).
		Save(ctx)
	if err != nil {
		return gerror.Wrap(err, "update task status fail")
	}
	return nil
}

func (f *filemgrRepo) UpdateTaskStatus(ctx context.Context, taskid, fileId string, status filemgr.Status) error {
	// 开始事务
	tx, err := f.db.Tx(ctx)
	if err != nil {
		return gerror.Wrap(err, "open tx fail")
	}
	_, err = tx.FileTransferTask.
		Update().
		Where(filetransfertask.TaskIDEQ(taskid)).
		SetStatus(status.Int()).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return gerror.Wrap(err, "update task status fail")
	}
	_, err = tx.SendFile.
		Update().
		Where(sendfile.FileIDEQ(fileId)).
		SetStatus(status.Int()).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return gerror.Wrap(err, "update sendfile status fail")
	}
	return tx.Commit()
}

func (f *filemgrRepo) UpdateRecvStatus(ctx context.Context, fileId string, status filemgr.Status) error {
	_, err := f.db.SendFile.
		Update().
		Where(sendfile.FileIDEQ(fileId)).
		SetStatus(status.Int()).
		Save(ctx)
	if err != nil {
		return gerror.Wrap(err, "update recvfile status fail")
	}
	return nil
}
