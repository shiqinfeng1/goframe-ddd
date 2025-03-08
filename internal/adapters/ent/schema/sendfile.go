package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SendFile holds the schema definition for the SendFile entity.
type SendFile struct {
	ent.Schema
}

// Fields of the SendFile.
func (SendFile) Fields() []ent.Field {
	return []ent.Field{
		field.String("task_id").NotEmpty().MaxLen(20),
		field.Text("task_name").NotEmpty(),
		field.Text("file_path").NotEmpty(),
		field.String("file_id").Unique().NotEmpty().MaxLen(20), // 标识文件的唯一id
		field.Int64("file_size").Positive(),
		field.Int("chunk_num_total").Positive(),
		field.Int("chunk_num_sended").Default(0),
		field.Int("status").Default(0),
		field.String("elapsed").Default("").MaxLen(16),
		field.String("speed").Default("").MaxLen(16),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the SendFile.
func (SendFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("send_chunks", SendChunk.Type),
	}
}
