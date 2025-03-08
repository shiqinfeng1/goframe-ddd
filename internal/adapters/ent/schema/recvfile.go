package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RecvFile holds the schema definition for the RecvFile entity.
type RecvFile struct {
	ent.Schema
}

// Fields of the RecvFile.
func (RecvFile) Fields() []ent.Field {
	return []ent.Field{
		field.String("task_id").Unique().NotEmpty().MaxLen(20),
		field.Text("task_name").NotEmpty(),
		field.Text("file_path_save").NotEmpty(),
		field.Text("file_path_origin").NotEmpty(),
		field.String("file_id").Unique().NotEmpty().MaxLen(20), // 标识文件的唯一id
		field.Int64("file_size").Positive(),
		field.Int("chunk_num_total").Positive(),
		field.Int("chunk_num_recved").Default(0),
		field.Int("status").Default(0),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the RecvFile.
func (RecvFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("recv_chunks", RecvChunk.Type),
	}
}
