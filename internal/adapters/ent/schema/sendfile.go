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
		field.Text("file_path").Optional().Nillable(),
		field.Text("file_name").Optional().Nillable(),
		field.String("fid").Unique(),
		field.Int64("file_size").Default(0),
		field.Int("chunk_num_total").Default(0),
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
