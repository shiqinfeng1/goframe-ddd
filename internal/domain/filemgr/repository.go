package filemgr

import "context"

// SendChunk is the model entity for the SendChunk schema.
type SendChunk struct {
	FileID      int   `json:"file_id,omitempty"` // 对应sendfile表的id
	ChunkIndex  int   `json:"chunk_index,omitempty"`
	ChunkOffset int64 `json:"chunk_offset,omitempty"`
	ChunkSize   int   `json:"chunk_size,omitempty"`
	Status      int   `json:"status,omitempty"`
}

// SendFile is the model entity for the SendFile schema.
type SendFile struct {
	ID             int    `json:"id,omitempty"`
	TaskID         string `json:"task_id,omitempty"`
	TaskName       string `json:"task_name,omitempty"`
	FilePath       string `json:"file_path,omitempty"`
	Fid            string `json:"fid,omitempty"`
	FileSize       int64  `json:"file_size,omitempty"`
	ChunkNumTotal  int    `json:"chunk_num_total,omitempty"`
	ChunkNumSended int    `json:"chunk_num_sended,omitempty"`
	Status         int    `json:"status,omitempty"`
	Elapsed        string `json:"elapsed,omitempty"`
	Speed          string `json:"speed,omitempty"`
}

type RecvChunk struct {
	FileID      int   `json:"file_id,omitempty"`
	ChunkIndex  int   `json:"chunk_index,omitempty"`
	ChunkOffset int64 `json:"chunk_offset,omitempty"`
	ChunkSize   int   `json:"chunk_size,omitempty"`
	Status      int   `json:"status,omitempty"`
}

// RecvFile is the model entity for the RecvFile schema.
type RecvFile struct {
	TaskID         string `json:"task_id,omitempty"`
	TaskName       string `json:"task_name,omitempty"`
	FilePathSave   string `json:"file_path_save,omitempty"`
	FilePathOrigin string `json:"file_path_origin,omitempty"`
	Fid            string `json:"fid,omitempty"`
	FileSize       int64  `json:"file_size,omitempty"`
	ChunkNumTotal  int    `json:"chunk_num_total,omitempty"`
	ChunkNumRecved int    `json:"chunk_num_recved,omitempty"`
	Status         int    `json:"status,omitempty"`
}

type Repository interface {
	GetSendFile(ctx context.Context, taskId, filePath string) (*SendFile, error)
	GetSendTask(ctx context.Context, taskId string) ([]*SendFile, error)
	SaveSendFile(ctx context.Context, sendFile *SendFile) (int, error)
	UpdateSendChunk(ctx context.Context, sendChunk *SendChunk) error
	GetRecvTask(ctx context.Context, taskId string) ([]*RecvFile, error)
	SaveRecvFile(ctx context.Context, recvile *RecvFile) error
	UpdateRecvChunk(ctx context.Context, recvChunk *RecvChunk) error
}
