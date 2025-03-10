package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SendFile holds the schema definition for the SendFile entity.
type FileTransferTask struct {
	ent.Schema
}

// Fields of the SendFile.
func (FileTransferTask) Fields() []ent.Field {
	return []ent.Field{
		field.String("task_id").Unique().NotEmpty().MaxLen(20),
		field.Text("task_name").NotEmpty(),
		field.Text("node_id").NotEmpty(),
		field.Int("status").Default(0),
		field.String("elapsed").Default("").MaxLen(16),
		field.String("speed").Default("").MaxLen(16),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Indexes 定义用户实体的索引
func (FileTransferTask) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("task_id").Unique(),
	}
}
