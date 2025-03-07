package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RecvChunk holds the schema definition for the RecvChunk entity.
type RecvChunk struct {
	ent.Schema
}

// Fields of the RecvChunk.
func (RecvChunk) Fields() []ent.Field {
	return []ent.Field{
		field.Int("file_id").Unique(),
		field.Int("chunk_index").Default(0),
		field.Int64("chunk_offset").Default(0),
		field.Int("chunk_size").Default(0),
		field.Int("status").Default(0),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the RecvChunk.
func (RecvChunk) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("recv_file", RecvFile.Type).
			Ref("recv_chunks").
			Field("file_id"). // 和recvfile表的id字段关联
			// 外键关联
			Unique().
			Required(),
	}
}
