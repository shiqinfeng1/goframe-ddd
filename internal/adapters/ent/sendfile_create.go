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

// SendFileCreate is the builder for creating a SendFile entity.
type SendFileCreate struct {
	config
	mutation *SendFileMutation
	hooks    []Hook
}

// SetTaskID sets the "task_id" field.
func (sfc *SendFileCreate) SetTaskID(s string) *SendFileCreate {
	sfc.mutation.SetTaskID(s)
	return sfc
}

// SetFilePath sets the "file_path" field.
func (sfc *SendFileCreate) SetFilePath(s string) *SendFileCreate {
	sfc.mutation.SetFilePath(s)
	return sfc
}

// SetFileID sets the "file_id" field.
func (sfc *SendFileCreate) SetFileID(s string) *SendFileCreate {
	sfc.mutation.SetFileID(s)
	return sfc
}

// SetFileSize sets the "file_size" field.
func (sfc *SendFileCreate) SetFileSize(i int64) *SendFileCreate {
	sfc.mutation.SetFileSize(i)
	return sfc
}

// SetChunkNumTotal sets the "chunk_num_total" field.
func (sfc *SendFileCreate) SetChunkNumTotal(i int) *SendFileCreate {
	sfc.mutation.SetChunkNumTotal(i)
	return sfc
}

// SetChunkNumSended sets the "chunk_num_sended" field.
func (sfc *SendFileCreate) SetChunkNumSended(i int) *SendFileCreate {
	sfc.mutation.SetChunkNumSended(i)
	return sfc
}

// SetNillableChunkNumSended sets the "chunk_num_sended" field if the given value is not nil.
func (sfc *SendFileCreate) SetNillableChunkNumSended(i *int) *SendFileCreate {
	if i != nil {
		sfc.SetChunkNumSended(*i)
	}
	return sfc
}

// SetStatus sets the "status" field.
func (sfc *SendFileCreate) SetStatus(i int) *SendFileCreate {
	sfc.mutation.SetStatus(i)
	return sfc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (sfc *SendFileCreate) SetNillableStatus(i *int) *SendFileCreate {
	if i != nil {
		sfc.SetStatus(*i)
	}
	return sfc
}

// SetUpdatedAt sets the "updated_at" field.
func (sfc *SendFileCreate) SetUpdatedAt(t time.Time) *SendFileCreate {
	sfc.mutation.SetUpdatedAt(t)
	return sfc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (sfc *SendFileCreate) SetNillableUpdatedAt(t *time.Time) *SendFileCreate {
	if t != nil {
		sfc.SetUpdatedAt(*t)
	}
	return sfc
}

// SetCreatedAt sets the "created_at" field.
func (sfc *SendFileCreate) SetCreatedAt(t time.Time) *SendFileCreate {
	sfc.mutation.SetCreatedAt(t)
	return sfc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (sfc *SendFileCreate) SetNillableCreatedAt(t *time.Time) *SendFileCreate {
	if t != nil {
		sfc.SetCreatedAt(*t)
	}
	return sfc
}

// AddSendChunkIDs adds the "send_chunks" edge to the SendChunk entity by IDs.
func (sfc *SendFileCreate) AddSendChunkIDs(ids ...int) *SendFileCreate {
	sfc.mutation.AddSendChunkIDs(ids...)
	return sfc
}

// AddSendChunks adds the "send_chunks" edges to the SendChunk entity.
func (sfc *SendFileCreate) AddSendChunks(s ...*SendChunk) *SendFileCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sfc.AddSendChunkIDs(ids...)
}

// Mutation returns the SendFileMutation object of the builder.
func (sfc *SendFileCreate) Mutation() *SendFileMutation {
	return sfc.mutation
}

// Save creates the SendFile in the database.
func (sfc *SendFileCreate) Save(ctx context.Context) (*SendFile, error) {
	sfc.defaults()
	return withHooks(ctx, sfc.sqlSave, sfc.mutation, sfc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sfc *SendFileCreate) SaveX(ctx context.Context) *SendFile {
	v, err := sfc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sfc *SendFileCreate) Exec(ctx context.Context) error {
	_, err := sfc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sfc *SendFileCreate) ExecX(ctx context.Context) {
	if err := sfc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sfc *SendFileCreate) defaults() {
	if _, ok := sfc.mutation.ChunkNumSended(); !ok {
		v := sendfile.DefaultChunkNumSended
		sfc.mutation.SetChunkNumSended(v)
	}
	if _, ok := sfc.mutation.Status(); !ok {
		v := sendfile.DefaultStatus
		sfc.mutation.SetStatus(v)
	}
	if _, ok := sfc.mutation.UpdatedAt(); !ok {
		v := sendfile.DefaultUpdatedAt()
		sfc.mutation.SetUpdatedAt(v)
	}
	if _, ok := sfc.mutation.CreatedAt(); !ok {
		v := sendfile.DefaultCreatedAt()
		sfc.mutation.SetCreatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sfc *SendFileCreate) check() error {
	if _, ok := sfc.mutation.TaskID(); !ok {
		return &ValidationError{Name: "task_id", err: errors.New(`ent: missing required field "SendFile.task_id"`)}
	}
	if v, ok := sfc.mutation.TaskID(); ok {
		if err := sendfile.TaskIDValidator(v); err != nil {
			return &ValidationError{Name: "task_id", err: fmt.Errorf(`ent: validator failed for field "SendFile.task_id": %w`, err)}
		}
	}
	if _, ok := sfc.mutation.FilePath(); !ok {
		return &ValidationError{Name: "file_path", err: errors.New(`ent: missing required field "SendFile.file_path"`)}
	}
	if v, ok := sfc.mutation.FilePath(); ok {
		if err := sendfile.FilePathValidator(v); err != nil {
			return &ValidationError{Name: "file_path", err: fmt.Errorf(`ent: validator failed for field "SendFile.file_path": %w`, err)}
		}
	}
	if _, ok := sfc.mutation.FileID(); !ok {
		return &ValidationError{Name: "file_id", err: errors.New(`ent: missing required field "SendFile.file_id"`)}
	}
	if v, ok := sfc.mutation.FileID(); ok {
		if err := sendfile.FileIDValidator(v); err != nil {
			return &ValidationError{Name: "file_id", err: fmt.Errorf(`ent: validator failed for field "SendFile.file_id": %w`, err)}
		}
	}
	if _, ok := sfc.mutation.FileSize(); !ok {
		return &ValidationError{Name: "file_size", err: errors.New(`ent: missing required field "SendFile.file_size"`)}
	}
	if v, ok := sfc.mutation.FileSize(); ok {
		if err := sendfile.FileSizeValidator(v); err != nil {
			return &ValidationError{Name: "file_size", err: fmt.Errorf(`ent: validator failed for field "SendFile.file_size": %w`, err)}
		}
	}
	if _, ok := sfc.mutation.ChunkNumTotal(); !ok {
		return &ValidationError{Name: "chunk_num_total", err: errors.New(`ent: missing required field "SendFile.chunk_num_total"`)}
	}
	if v, ok := sfc.mutation.ChunkNumTotal(); ok {
		if err := sendfile.ChunkNumTotalValidator(v); err != nil {
			return &ValidationError{Name: "chunk_num_total", err: fmt.Errorf(`ent: validator failed for field "SendFile.chunk_num_total": %w`, err)}
		}
	}
	if _, ok := sfc.mutation.ChunkNumSended(); !ok {
		return &ValidationError{Name: "chunk_num_sended", err: errors.New(`ent: missing required field "SendFile.chunk_num_sended"`)}
	}
	if _, ok := sfc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "SendFile.status"`)}
	}
	if _, ok := sfc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "SendFile.updated_at"`)}
	}
	if _, ok := sfc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "SendFile.created_at"`)}
	}
	return nil
}

func (sfc *SendFileCreate) sqlSave(ctx context.Context) (*SendFile, error) {
	if err := sfc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sfc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sfc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sfc.mutation.id = &_node.ID
	sfc.mutation.done = true
	return _node, nil
}

func (sfc *SendFileCreate) createSpec() (*SendFile, *sqlgraph.CreateSpec) {
	var (
		_node = &SendFile{config: sfc.config}
		_spec = sqlgraph.NewCreateSpec(sendfile.Table, sqlgraph.NewFieldSpec(sendfile.FieldID, field.TypeInt))
	)
	if value, ok := sfc.mutation.TaskID(); ok {
		_spec.SetField(sendfile.FieldTaskID, field.TypeString, value)
		_node.TaskID = value
	}
	if value, ok := sfc.mutation.FilePath(); ok {
		_spec.SetField(sendfile.FieldFilePath, field.TypeString, value)
		_node.FilePath = value
	}
	if value, ok := sfc.mutation.FileID(); ok {
		_spec.SetField(sendfile.FieldFileID, field.TypeString, value)
		_node.FileID = value
	}
	if value, ok := sfc.mutation.FileSize(); ok {
		_spec.SetField(sendfile.FieldFileSize, field.TypeInt64, value)
		_node.FileSize = value
	}
	if value, ok := sfc.mutation.ChunkNumTotal(); ok {
		_spec.SetField(sendfile.FieldChunkNumTotal, field.TypeInt, value)
		_node.ChunkNumTotal = value
	}
	if value, ok := sfc.mutation.ChunkNumSended(); ok {
		_spec.SetField(sendfile.FieldChunkNumSended, field.TypeInt, value)
		_node.ChunkNumSended = value
	}
	if value, ok := sfc.mutation.Status(); ok {
		_spec.SetField(sendfile.FieldStatus, field.TypeInt, value)
		_node.Status = value
	}
	if value, ok := sfc.mutation.UpdatedAt(); ok {
		_spec.SetField(sendfile.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := sfc.mutation.CreatedAt(); ok {
		_spec.SetField(sendfile.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := sfc.mutation.SendChunksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   sendfile.SendChunksTable,
			Columns: []string{sendfile.SendChunksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(sendchunk.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SendFileCreateBulk is the builder for creating many SendFile entities in bulk.
type SendFileCreateBulk struct {
	config
	err      error
	builders []*SendFileCreate
}

// Save creates the SendFile entities in the database.
func (sfcb *SendFileCreateBulk) Save(ctx context.Context) ([]*SendFile, error) {
	if sfcb.err != nil {
		return nil, sfcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(sfcb.builders))
	nodes := make([]*SendFile, len(sfcb.builders))
	mutators := make([]Mutator, len(sfcb.builders))
	for i := range sfcb.builders {
		func(i int, root context.Context) {
			builder := sfcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SendFileMutation)
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
					_, err = mutators[i+1].Mutate(root, sfcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, sfcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, sfcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (sfcb *SendFileCreateBulk) SaveX(ctx context.Context) []*SendFile {
	v, err := sfcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sfcb *SendFileCreateBulk) Exec(ctx context.Context) error {
	_, err := sfcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sfcb *SendFileCreateBulk) ExecX(ctx context.Context) {
	if err := sfcb.Exec(ctx); err != nil {
		panic(err)
	}
}
