// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/sendchunk"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/sendfile"
)

// SendChunkCreate is the builder for creating a SendChunk entity.
type SendChunkCreate struct {
	config
	mutation *SendChunkMutation
	hooks    []Hook
}

// SetSendfileID sets the "sendfile_id" field.
func (scc *SendChunkCreate) SetSendfileID(i int) *SendChunkCreate {
	scc.mutation.SetSendfileID(i)
	return scc
}

// SetChunkIndex sets the "chunk_index" field.
func (scc *SendChunkCreate) SetChunkIndex(i int) *SendChunkCreate {
	scc.mutation.SetChunkIndex(i)
	return scc
}

// SetNillableChunkIndex sets the "chunk_index" field if the given value is not nil.
func (scc *SendChunkCreate) SetNillableChunkIndex(i *int) *SendChunkCreate {
	if i != nil {
		scc.SetChunkIndex(*i)
	}
	return scc
}

// SetChunkOffset sets the "chunk_offset" field.
func (scc *SendChunkCreate) SetChunkOffset(i int64) *SendChunkCreate {
	scc.mutation.SetChunkOffset(i)
	return scc
}

// SetNillableChunkOffset sets the "chunk_offset" field if the given value is not nil.
func (scc *SendChunkCreate) SetNillableChunkOffset(i *int64) *SendChunkCreate {
	if i != nil {
		scc.SetChunkOffset(*i)
	}
	return scc
}

// SetChunkSize sets the "chunk_size" field.
func (scc *SendChunkCreate) SetChunkSize(i int) *SendChunkCreate {
	scc.mutation.SetChunkSize(i)
	return scc
}

// SetNillableChunkSize sets the "chunk_size" field if the given value is not nil.
func (scc *SendChunkCreate) SetNillableChunkSize(i *int) *SendChunkCreate {
	if i != nil {
		scc.SetChunkSize(*i)
	}
	return scc
}

// SetUpdatedAt sets the "updated_at" field.
func (scc *SendChunkCreate) SetUpdatedAt(t time.Time) *SendChunkCreate {
	scc.mutation.SetUpdatedAt(t)
	return scc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (scc *SendChunkCreate) SetNillableUpdatedAt(t *time.Time) *SendChunkCreate {
	if t != nil {
		scc.SetUpdatedAt(*t)
	}
	return scc
}

// SetCreatedAt sets the "created_at" field.
func (scc *SendChunkCreate) SetCreatedAt(t time.Time) *SendChunkCreate {
	scc.mutation.SetCreatedAt(t)
	return scc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (scc *SendChunkCreate) SetNillableCreatedAt(t *time.Time) *SendChunkCreate {
	if t != nil {
		scc.SetCreatedAt(*t)
	}
	return scc
}

// SetSendFileID sets the "send_file" edge to the SendFile entity by ID.
func (scc *SendChunkCreate) SetSendFileID(id int) *SendChunkCreate {
	scc.mutation.SetSendFileID(id)
	return scc
}

// SetSendFile sets the "send_file" edge to the SendFile entity.
func (scc *SendChunkCreate) SetSendFile(s *SendFile) *SendChunkCreate {
	return scc.SetSendFileID(s.ID)
}

// Mutation returns the SendChunkMutation object of the builder.
func (scc *SendChunkCreate) Mutation() *SendChunkMutation {
	return scc.mutation
}

// Save creates the SendChunk in the database.
func (scc *SendChunkCreate) Save(ctx context.Context) (*SendChunk, error) {
	scc.defaults()
	return withHooks(ctx, scc.sqlSave, scc.mutation, scc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (scc *SendChunkCreate) SaveX(ctx context.Context) *SendChunk {
	v, err := scc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scc *SendChunkCreate) Exec(ctx context.Context) error {
	_, err := scc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scc *SendChunkCreate) ExecX(ctx context.Context) {
	if err := scc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (scc *SendChunkCreate) defaults() {
	if _, ok := scc.mutation.ChunkIndex(); !ok {
		v := sendchunk.DefaultChunkIndex
		scc.mutation.SetChunkIndex(v)
	}
	if _, ok := scc.mutation.ChunkOffset(); !ok {
		v := sendchunk.DefaultChunkOffset
		scc.mutation.SetChunkOffset(v)
	}
	if _, ok := scc.mutation.ChunkSize(); !ok {
		v := sendchunk.DefaultChunkSize
		scc.mutation.SetChunkSize(v)
	}
	if _, ok := scc.mutation.UpdatedAt(); !ok {
		v := sendchunk.DefaultUpdatedAt()
		scc.mutation.SetUpdatedAt(v)
	}
	if _, ok := scc.mutation.CreatedAt(); !ok {
		v := sendchunk.DefaultCreatedAt()
		scc.mutation.SetCreatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (scc *SendChunkCreate) check() error {
	if _, ok := scc.mutation.SendfileID(); !ok {
		return &ValidationError{Name: "sendfile_id", err: errors.New(`ent: missing required field "SendChunk.sendfile_id"`)}
	}
	if _, ok := scc.mutation.ChunkIndex(); !ok {
		return &ValidationError{Name: "chunk_index", err: errors.New(`ent: missing required field "SendChunk.chunk_index"`)}
	}
	if _, ok := scc.mutation.ChunkOffset(); !ok {
		return &ValidationError{Name: "chunk_offset", err: errors.New(`ent: missing required field "SendChunk.chunk_offset"`)}
	}
	if _, ok := scc.mutation.ChunkSize(); !ok {
		return &ValidationError{Name: "chunk_size", err: errors.New(`ent: missing required field "SendChunk.chunk_size"`)}
	}
	if _, ok := scc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "SendChunk.updated_at"`)}
	}
	if _, ok := scc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "SendChunk.created_at"`)}
	}
	if len(scc.mutation.SendFileIDs()) == 0 {
		return &ValidationError{Name: "send_file", err: errors.New(`ent: missing required edge "SendChunk.send_file"`)}
	}
	return nil
}

func (scc *SendChunkCreate) sqlSave(ctx context.Context) (*SendChunk, error) {
	if err := scc.check(); err != nil {
		return nil, err
	}
	_node, _spec := scc.createSpec()
	if err := sqlgraph.CreateNode(ctx, scc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	scc.mutation.id = &_node.ID
	scc.mutation.done = true
	return _node, nil
}

func (scc *SendChunkCreate) createSpec() (*SendChunk, *sqlgraph.CreateSpec) {
	var (
		_node = &SendChunk{config: scc.config}
		_spec = sqlgraph.NewCreateSpec(sendchunk.Table, sqlgraph.NewFieldSpec(sendchunk.FieldID, field.TypeInt))
	)
	if value, ok := scc.mutation.ChunkIndex(); ok {
		_spec.SetField(sendchunk.FieldChunkIndex, field.TypeInt, value)
		_node.ChunkIndex = value
	}
	if value, ok := scc.mutation.ChunkOffset(); ok {
		_spec.SetField(sendchunk.FieldChunkOffset, field.TypeInt64, value)
		_node.ChunkOffset = value
	}
	if value, ok := scc.mutation.ChunkSize(); ok {
		_spec.SetField(sendchunk.FieldChunkSize, field.TypeInt, value)
		_node.ChunkSize = value
	}
	if value, ok := scc.mutation.UpdatedAt(); ok {
		_spec.SetField(sendchunk.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := scc.mutation.CreatedAt(); ok {
		_spec.SetField(sendchunk.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := scc.mutation.SendFileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   sendchunk.SendFileTable,
			Columns: []string{sendchunk.SendFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(sendfile.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.SendfileID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SendChunkCreateBulk is the builder for creating many SendChunk entities in bulk.
type SendChunkCreateBulk struct {
	config
	err      error
	builders []*SendChunkCreate
}

// Save creates the SendChunk entities in the database.
func (sccb *SendChunkCreateBulk) Save(ctx context.Context) ([]*SendChunk, error) {
	if sccb.err != nil {
		return nil, sccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(sccb.builders))
	nodes := make([]*SendChunk, len(sccb.builders))
	mutators := make([]Mutator, len(sccb.builders))
	for i := range sccb.builders {
		func(i int, root context.Context) {
			builder := sccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SendChunkMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, sccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, sccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, sccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (sccb *SendChunkCreateBulk) SaveX(ctx context.Context) []*SendChunk {
	v, err := sccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sccb *SendChunkCreateBulk) Exec(ctx context.Context) error {
	_, err := sccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sccb *SendChunkCreateBulk) ExecX(ctx context.Context) {
	if err := sccb.Exec(ctx); err != nil {
		panic(err)
	}
}
