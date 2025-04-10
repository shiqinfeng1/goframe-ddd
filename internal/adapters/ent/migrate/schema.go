// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// FileTransferTasksColumns holds the columns for the "file_transfer_tasks" table.
	FileTransferTasksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "task_id", Type: field.TypeString, Unique: true, Size: 20},
		{Name: "task_name", Type: field.TypeString, Size: 2147483647},
		{Name: "node_id", Type: field.TypeString, Size: 2147483647},
		{Name: "status", Type: field.TypeInt, Default: 0},
		{Name: "elapsed", Type: field.TypeString, Size: 16, Default: ""},
		{Name: "speed", Type: field.TypeString, Size: 16, Default: ""},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "created_at", Type: field.TypeTime},
	}
	// FileTransferTasksTable holds the schema information for the "file_transfer_tasks" table.
	FileTransferTasksTable = &schema.Table{
		Name:       "file_transfer_tasks",
		Columns:    FileTransferTasksColumns,
		PrimaryKey: []*schema.Column{FileTransferTasksColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "filetransfertask_task_id",
				Unique:  true,
				Columns: []*schema.Column{FileTransferTasksColumns[1]},
			},
		},
	}
	// RecvChunksColumns holds the columns for the "recv_chunks" table.
	RecvChunksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "chunk_index", Type: field.TypeInt, Default: 0},
		{Name: "chunk_offset", Type: field.TypeInt64, Default: 0},
		{Name: "chunk_size", Type: field.TypeInt, Default: 0},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "recvfile_id", Type: field.TypeInt},
	}
	// RecvChunksTable holds the schema information for the "recv_chunks" table.
	RecvChunksTable = &schema.Table{
		Name:       "recv_chunks",
		Columns:    RecvChunksColumns,
		PrimaryKey: []*schema.Column{RecvChunksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "recv_chunks_recv_files_recv_chunks",
				Columns:    []*schema.Column{RecvChunksColumns[6]},
				RefColumns: []*schema.Column{RecvFilesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// RecvFilesColumns holds the columns for the "recv_files" table.
	RecvFilesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "task_id", Type: field.TypeString, Unique: true, Size: 20},
		{Name: "task_name", Type: field.TypeString, Size: 2147483647},
		{Name: "file_path_save", Type: field.TypeString, Size: 2147483647},
		{Name: "file_path_origin", Type: field.TypeString, Size: 2147483647},
		{Name: "file_id", Type: field.TypeString, Unique: true, Size: 20},
		{Name: "file_size", Type: field.TypeInt64},
		{Name: "chunk_num_total", Type: field.TypeInt},
		{Name: "chunk_num_recved", Type: field.TypeInt, Default: 0},
		{Name: "status", Type: field.TypeInt, Default: 0},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "created_at", Type: field.TypeTime},
	}
	// RecvFilesTable holds the schema information for the "recv_files" table.
	RecvFilesTable = &schema.Table{
		Name:       "recv_files",
		Columns:    RecvFilesColumns,
		PrimaryKey: []*schema.Column{RecvFilesColumns[0]},
	}
	// SendChunksColumns holds the columns for the "send_chunks" table.
	SendChunksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "chunk_index", Type: field.TypeInt, Default: 0},
		{Name: "chunk_offset", Type: field.TypeInt64, Default: 0},
		{Name: "chunk_size", Type: field.TypeInt, Default: 0},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "sendfile_id", Type: field.TypeInt},
	}
	// SendChunksTable holds the schema information for the "send_chunks" table.
	SendChunksTable = &schema.Table{
		Name:       "send_chunks",
		Columns:    SendChunksColumns,
		PrimaryKey: []*schema.Column{SendChunksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "send_chunks_send_files_send_chunks",
				Columns:    []*schema.Column{SendChunksColumns[6]},
				RefColumns: []*schema.Column{SendFilesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// SendFilesColumns holds the columns for the "send_files" table.
	SendFilesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "task_id", Type: field.TypeString, Unique: true, Size: 20},
		{Name: "file_path", Type: field.TypeString, Size: 2147483647},
		{Name: "file_id", Type: field.TypeString, Unique: true, Size: 20},
		{Name: "file_size", Type: field.TypeInt64},
		{Name: "chunk_num_total", Type: field.TypeInt},
		{Name: "chunk_num_sended", Type: field.TypeInt, Default: 0},
		{Name: "status", Type: field.TypeInt, Default: 0},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "created_at", Type: field.TypeTime},
	}
	// SendFilesTable holds the schema information for the "send_files" table.
	SendFilesTable = &schema.Table{
		Name:       "send_files",
		Columns:    SendFilesColumns,
		PrimaryKey: []*schema.Column{SendFilesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "sendfile_task_id_file_path",
				Unique:  true,
				Columns: []*schema.Column{SendFilesColumns[1], SendFilesColumns[2]},
			},
			{
				Name:    "sendfile_file_id",
				Unique:  true,
				Columns: []*schema.Column{SendFilesColumns[3]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		FileTransferTasksTable,
		RecvChunksTable,
		RecvFilesTable,
		SendChunksTable,
		SendFilesTable,
	}
)

func init() {
	RecvChunksTable.ForeignKeys[0].RefTable = RecvFilesTable
	SendChunksTable.ForeignKeys[0].RefTable = SendFilesTable
}
