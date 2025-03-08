package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SendChunk holds the schema definition for the SendChunk entity.
type SendChunk struct {
	ent.Schema
}

// Fields of the SendChunk.
func (SendChunk) Fields() []ent.Field {
	return []ent.Field{
		field.Int("sendfile_id").Unique(),
		field.Int("chunk_index").Default(0),
		field.Int64("chunk_offset").Default(0),
		field.Int("chunk_size").Default(0),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the SendChunk.
func (SendChunk) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("send_file", SendFile.Type).
			Ref("send_chunks").
			Field("sendfile_id"). // 和sendfile表的id字段关联
			// 外键关联
			Unique().
			Required(),
	}
}
