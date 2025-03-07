// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/predicate"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/sendchunk"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/sendfile"
)

// SendFileUpdate is the builder for updating SendFile entities.
type SendFileUpdate struct {
	config
	hooks    []Hook
	mutation *SendFileMutation
}

// Where appends a list predicates to the SendFileUpdate builder.
func (sfu *SendFileUpdate) Where(ps ...predicate.SendFile) *SendFileUpdate {
	sfu.mutation.Where(ps...)
	return sfu
}

// SetTaskID sets the "task_id" field.
func (sfu *SendFileUpdate) SetTaskID(s string) *SendFileUpdate {
	sfu.mutation.SetTaskID(s)
	return sfu
}

// SetNillableTaskID sets the "task_id" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableTaskID(s *string) *SendFileUpdate {
	if s != nil {
		sfu.SetTaskID(*s)
	}
	return sfu
}

// SetTaskName sets the "task_name" field.
func (sfu *SendFileUpdate) SetTaskName(s string) *SendFileUpdate {
	sfu.mutation.SetTaskName(s)
	return sfu
}

// SetNillableTaskName sets the "task_name" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableTaskName(s *string) *SendFileUpdate {
	if s != nil {
		sfu.SetTaskName(*s)
	}
	return sfu
}

// SetFilePath sets the "file_path" field.
func (sfu *SendFileUpdate) SetFilePath(s string) *SendFileUpdate {
	sfu.mutation.SetFilePath(s)
	return sfu
}

// SetNillableFilePath sets the "file_path" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableFilePath(s *string) *SendFileUpdate {
	if s != nil {
		sfu.SetFilePath(*s)
	}
	return sfu
}

// SetFid sets the "fid" field.
func (sfu *SendFileUpdate) SetFid(s string) *SendFileUpdate {
	sfu.mutation.SetFid(s)
	return sfu
}

// SetNillableFid sets the "fid" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableFid(s *string) *SendFileUpdate {
	if s != nil {
		sfu.SetFid(*s)
	}
	return sfu
}

// SetFileSize sets the "file_size" field.
func (sfu *SendFileUpdate) SetFileSize(i int64) *SendFileUpdate {
	sfu.mutation.ResetFileSize()
	sfu.mutation.SetFileSize(i)
	return sfu
}

// SetNillableFileSize sets the "file_size" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableFileSize(i *int64) *SendFileUpdate {
	if i != nil {
		sfu.SetFileSize(*i)
	}
	return sfu
}

// AddFileSize adds i to the "file_size" field.
func (sfu *SendFileUpdate) AddFileSize(i int64) *SendFileUpdate {
	sfu.mutation.AddFileSize(i)
	return sfu
}

// SetChunkNumTotal sets the "chunk_num_total" field.
func (sfu *SendFileUpdate) SetChunkNumTotal(i int) *SendFileUpdate {
	sfu.mutation.ResetChunkNumTotal()
	sfu.mutation.SetChunkNumTotal(i)
	return sfu
}

// SetNillableChunkNumTotal sets the "chunk_num_total" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableChunkNumTotal(i *int) *SendFileUpdate {
	if i != nil {
		sfu.SetChunkNumTotal(*i)
	}
	return sfu
}

// AddChunkNumTotal adds i to the "chunk_num_total" field.
func (sfu *SendFileUpdate) AddChunkNumTotal(i int) *SendFileUpdate {
	sfu.mutation.AddChunkNumTotal(i)
	return sfu
}

// SetChunkNumSended sets the "chunk_num_sended" field.
func (sfu *SendFileUpdate) SetChunkNumSended(i int) *SendFileUpdate {
	sfu.mutation.ResetChunkNumSended()
	sfu.mutation.SetChunkNumSended(i)
	return sfu
}

// SetNillableChunkNumSended sets the "chunk_num_sended" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableChunkNumSended(i *int) *SendFileUpdate {
	if i != nil {
		sfu.SetChunkNumSended(*i)
	}
	return sfu
}

// AddChunkNumSended adds i to the "chunk_num_sended" field.
func (sfu *SendFileUpdate) AddChunkNumSended(i int) *SendFileUpdate {
	sfu.mutation.AddChunkNumSended(i)
	return sfu
}

// SetStatus sets the "status" field.
func (sfu *SendFileUpdate) SetStatus(i int) *SendFileUpdate {
	sfu.mutation.ResetStatus()
	sfu.mutation.SetStatus(i)
	return sfu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableStatus(i *int) *SendFileUpdate {
	if i != nil {
		sfu.SetStatus(*i)
	}
	return sfu
}

// AddStatus adds i to the "status" field.
func (sfu *SendFileUpdate) AddStatus(i int) *SendFileUpdate {
	sfu.mutation.AddStatus(i)
	return sfu
}

// SetElapsed sets the "elapsed" field.
func (sfu *SendFileUpdate) SetElapsed(s string) *SendFileUpdate {
	sfu.mutation.SetElapsed(s)
	return sfu
}

// SetNillableElapsed sets the "elapsed" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableElapsed(s *string) *SendFileUpdate {
	if s != nil {
		sfu.SetElapsed(*s)
	}
	return sfu
}

// SetSpeed sets the "speed" field.
func (sfu *SendFileUpdate) SetSpeed(s string) *SendFileUpdate {
	sfu.mutation.SetSpeed(s)
	return sfu
}

// SetNillableSpeed sets the "speed" field if the given value is not nil.
func (sfu *SendFileUpdate) SetNillableSpeed(s *string) *SendFileUpdate {
	if s != nil {
		sfu.SetSpeed(*s)
	}
	return sfu
}

// SetUpdatedAt sets the "updated_at" field.
func (sfu *SendFileUpdate) SetUpdatedAt(t time.Time) *SendFileUpdate {
	sfu.mutation.SetUpdatedAt(t)
	return sfu
}

// AddSendChunkIDs adds the "send_chunks" edge to the SendChunk entity by IDs.
func (sfu *SendFileUpdate) AddSendChunkIDs(ids ...int) *SendFileUpdate {
	sfu.mutation.AddSendChunkIDs(ids...)
	return sfu
}

// AddSendChunks adds the "send_chunks" edges to the SendChunk entity.
func (sfu *SendFileUpdate) AddSendChunks(s ...*SendChunk) *SendFileUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sfu.AddSendChunkIDs(ids...)
}

// Mutation returns the SendFileMutation object of the builder.
func (sfu *SendFileUpdate) Mutation() *SendFileMutation {
	return sfu.mutation
}

// ClearSendChunks clears all "send_chunks" edges to the SendChunk entity.
func (sfu *SendFileUpdate) ClearSendChunks() *SendFileUpdate {
	sfu.mutation.ClearSendChunks()
	return sfu
}

// RemoveSendChunkIDs removes the "send_chunks" edge to SendChunk entities by IDs.
func (sfu *SendFileUpdate) RemoveSendChunkIDs(ids ...int) *SendFileUpdate {
	sfu.mutation.RemoveSendChunkIDs(ids...)
	return sfu
}

// RemoveSendChunks removes "send_chunks" edges to SendChunk entities.
func (sfu *SendFileUpdate) RemoveSendChunks(s ...*SendChunk) *SendFileUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sfu.RemoveSendChunkIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (sfu *SendFileUpdate) Save(ctx context.Context) (int, error) {
	sfu.defaults()
	return withHooks(ctx, sfu.sqlSave, sfu.mutation, sfu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (sfu *SendFileUpdate) SaveX(ctx context.Context) int {
	affected, err := sfu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (sfu *SendFileUpdate) Exec(ctx context.Context) error {
	_, err := sfu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sfu *SendFileUpdate) ExecX(ctx context.Context) {
	if err := sfu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sfu *SendFileUpdate) defaults() {
	if _, ok := sfu.mutation.UpdatedAt(); !ok {
		v := sendfile.UpdateDefaultUpdatedAt()
		sfu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sfu *SendFileUpdate) check() error {
	if v, ok := sfu.mutation.TaskID(); ok {
		if err := sendfile.TaskIDValidator(v); err != nil {
			return &ValidationError{Name: "task_id", err: fmt.Errorf(`ent: validator failed for field "SendFile.task_id": %w`, err)}
		}
	}
	if v, ok := sfu.mutation.TaskName(); ok {
		if err := sendfile.TaskNameValidator(v); err != nil {
			return &ValidationError{Name: "task_name", err: fmt.Errorf(`ent: validator failed for field "SendFile.task_name": %w`, err)}
		}
	}
	if v, ok := sfu.mutation.FilePath(); ok {
		if err := sendfile.FilePathValidator(v); err != nil {
			return &ValidationError{Name: "file_path", err: fmt.Errorf(`ent: validator failed for field "SendFile.file_path": %w`, err)}
		}
	}
	if v, ok := sfu.mutation.Fid(); ok {
		if err := sendfile.FidValidator(v); err != nil {
			return &ValidationError{Name: "fid", err: fmt.Errorf(`ent: validator failed for field "SendFile.fid": %w`, err)}
		}
	}
	if v, ok := sfu.mutation.Elapsed(); ok {
		if err := sendfile.ElapsedValidator(v); err != nil {
			return &ValidationError{Name: "elapsed", err: fmt.Errorf(`ent: validator failed for field "SendFile.elapsed": %w`, err)}
		}
	}
	if v, ok := sfu.mutation.Speed(); ok {
		if err := sendfile.SpeedValidator(v); err != nil {
			return &ValidationError{Name: "speed", err: fmt.Errorf(`ent: validator failed for field "SendFile.speed": %w`, err)}
		}
	}
	return nil
}

func (sfu *SendFileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := sfu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(sendfile.Table, sendfile.Columns, sqlgraph.NewFieldSpec(sendfile.FieldID, field.TypeInt))
	if ps := sfu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sfu.mutation.TaskID(); ok {
		_spec.SetField(sendfile.FieldTaskID, field.TypeString, value)
	}
	if value, ok := sfu.mutation.TaskName(); ok {
		_spec.SetField(sendfile.FieldTaskName, field.TypeString, value)
	}
	if value, ok := sfu.mutation.FilePath(); ok {
		_spec.SetField(sendfile.FieldFilePath, field.TypeString, value)
	}
	if value, ok := sfu.mutation.Fid(); ok {
		_spec.SetField(sendfile.FieldFid, field.TypeString, value)
	}
	if value, ok := sfu.mutation.FileSize(); ok {
		_spec.SetField(sendfile.FieldFileSize, field.TypeInt64, value)
	}
	if value, ok := sfu.mutation.AddedFileSize(); ok {
		_spec.AddField(sendfile.FieldFileSize, field.TypeInt64, value)
	}
	if value, ok := sfu.mutation.ChunkNumTotal(); ok {
		_spec.SetField(sendfile.FieldChunkNumTotal, field.TypeInt, value)
	}
	if value, ok := sfu.mutation.AddedChunkNumTotal(); ok {
		_spec.AddField(sendfile.FieldChunkNumTotal, field.TypeInt, value)
	}
	if value, ok := sfu.mutation.ChunkNumSended(); ok {
		_spec.SetField(sendfile.FieldChunkNumSended, field.TypeInt, value)
	}
	if value, ok := sfu.mutation.AddedChunkNumSended(); ok {
		_spec.AddField(sendfile.FieldChunkNumSended, field.TypeInt, value)
	}
	if value, ok := sfu.mutation.Status(); ok {
		_spec.SetField(sendfile.FieldStatus, field.TypeInt, value)
	}
	if value, ok := sfu.mutation.AddedStatus(); ok {
		_spec.AddField(sendfile.FieldStatus, field.TypeInt, value)
	}
	if value, ok := sfu.mutation.Elapsed(); ok {
		_spec.SetField(sendfile.FieldElapsed, field.TypeString, value)
	}
	if value, ok := sfu.mutation.Speed(); ok {
		_spec.SetField(sendfile.FieldSpeed, field.TypeString, value)
	}
	if value, ok := sfu.mutation.UpdatedAt(); ok {
		_spec.SetField(sendfile.FieldUpdatedAt, field.TypeTime, value)
	}
	if sfu.mutation.SendChunksCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := sfu.mutation.RemovedSendChunksIDs(); len(nodes) > 0 && !sfu.mutation.SendChunksCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := sfu.mutation.SendChunksIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, sfu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{sendfile.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	sfu.mutation.done = true
	return n, nil
}

// SendFileUpdateOne is the builder for updating a single SendFile entity.
type SendFileUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SendFileMutation
}

// SetTaskID sets the "task_id" field.
func (sfuo *SendFileUpdateOne) SetTaskID(s string) *SendFileUpdateOne {
	sfuo.mutation.SetTaskID(s)
	return sfuo
}

// SetNillableTaskID sets the "task_id" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableTaskID(s *string) *SendFileUpdateOne {
	if s != nil {
		sfuo.SetTaskID(*s)
	}
	return sfuo
}

// SetTaskName sets the "task_name" field.
func (sfuo *SendFileUpdateOne) SetTaskName(s string) *SendFileUpdateOne {
	sfuo.mutation.SetTaskName(s)
	return sfuo
}

// SetNillableTaskName sets the "task_name" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableTaskName(s *string) *SendFileUpdateOne {
	if s != nil {
		sfuo.SetTaskName(*s)
	}
	return sfuo
}

// SetFilePath sets the "file_path" field.
func (sfuo *SendFileUpdateOne) SetFilePath(s string) *SendFileUpdateOne {
	sfuo.mutation.SetFilePath(s)
	return sfuo
}

// SetNillableFilePath sets the "file_path" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableFilePath(s *string) *SendFileUpdateOne {
	if s != nil {
		sfuo.SetFilePath(*s)
	}
	return sfuo
}

// SetFid sets the "fid" field.
func (sfuo *SendFileUpdateOne) SetFid(s string) *SendFileUpdateOne {
	sfuo.mutation.SetFid(s)
	return sfuo
}

// SetNillableFid sets the "fid" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableFid(s *string) *SendFileUpdateOne {
	if s != nil {
		sfuo.SetFid(*s)
	}
	return sfuo
}

// SetFileSize sets the "file_size" field.
func (sfuo *SendFileUpdateOne) SetFileSize(i int64) *SendFileUpdateOne {
	sfuo.mutation.ResetFileSize()
	sfuo.mutation.SetFileSize(i)
	return sfuo
}

// SetNillableFileSize sets the "file_size" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableFileSize(i *int64) *SendFileUpdateOne {
	if i != nil {
		sfuo.SetFileSize(*i)
	}
	return sfuo
}

// AddFileSize adds i to the "file_size" field.
func (sfuo *SendFileUpdateOne) AddFileSize(i int64) *SendFileUpdateOne {
	sfuo.mutation.AddFileSize(i)
	return sfuo
}

// SetChunkNumTotal sets the "chunk_num_total" field.
func (sfuo *SendFileUpdateOne) SetChunkNumTotal(i int) *SendFileUpdateOne {
	sfuo.mutation.ResetChunkNumTotal()
	sfuo.mutation.SetChunkNumTotal(i)
	return sfuo
}

// SetNillableChunkNumTotal sets the "chunk_num_total" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableChunkNumTotal(i *int) *SendFileUpdateOne {
	if i != nil {
		sfuo.SetChunkNumTotal(*i)
	}
	return sfuo
}

// AddChunkNumTotal adds i to the "chunk_num_total" field.
func (sfuo *SendFileUpdateOne) AddChunkNumTotal(i int) *SendFileUpdateOne {
	sfuo.mutation.AddChunkNumTotal(i)
	return sfuo
}

// SetChunkNumSended sets the "chunk_num_sended" field.
func (sfuo *SendFileUpdateOne) SetChunkNumSended(i int) *SendFileUpdateOne {
	sfuo.mutation.ResetChunkNumSended()
	sfuo.mutation.SetChunkNumSended(i)
	return sfuo
}

// SetNillableChunkNumSended sets the "chunk_num_sended" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableChunkNumSended(i *int) *SendFileUpdateOne {
	if i != nil {
		sfuo.SetChunkNumSended(*i)
	}
	return sfuo
}

// AddChunkNumSended adds i to the "chunk_num_sended" field.
func (sfuo *SendFileUpdateOne) AddChunkNumSended(i int) *SendFileUpdateOne {
	sfuo.mutation.AddChunkNumSended(i)
	return sfuo
}

// SetStatus sets the "status" field.
func (sfuo *SendFileUpdateOne) SetStatus(i int) *SendFileUpdateOne {
	sfuo.mutation.ResetStatus()
	sfuo.mutation.SetStatus(i)
	return sfuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableStatus(i *int) *SendFileUpdateOne {
	if i != nil {
		sfuo.SetStatus(*i)
	}
	return sfuo
}

// AddStatus adds i to the "status" field.
func (sfuo *SendFileUpdateOne) AddStatus(i int) *SendFileUpdateOne {
	sfuo.mutation.AddStatus(i)
	return sfuo
}

// SetElapsed sets the "elapsed" field.
func (sfuo *SendFileUpdateOne) SetElapsed(s string) *SendFileUpdateOne {
	sfuo.mutation.SetElapsed(s)
	return sfuo
}

// SetNillableElapsed sets the "elapsed" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableElapsed(s *string) *SendFileUpdateOne {
	if s != nil {
		sfuo.SetElapsed(*s)
	}
	return sfuo
}

// SetSpeed sets the "speed" field.
func (sfuo *SendFileUpdateOne) SetSpeed(s string) *SendFileUpdateOne {
	sfuo.mutation.SetSpeed(s)
	return sfuo
}

// SetNillableSpeed sets the "speed" field if the given value is not nil.
func (sfuo *SendFileUpdateOne) SetNillableSpeed(s *string) *SendFileUpdateOne {
	if s != nil {
		sfuo.SetSpeed(*s)
	}
	return sfuo
}

// SetUpdatedAt sets the "updated_at" field.
func (sfuo *SendFileUpdateOne) SetUpdatedAt(t time.Time) *SendFileUpdateOne {
	sfuo.mutation.SetUpdatedAt(t)
	return sfuo
}

// AddSendChunkIDs adds the "send_chunks" edge to the SendChunk entity by IDs.
func (sfuo *SendFileUpdateOne) AddSendChunkIDs(ids ...int) *SendFileUpdateOne {
	sfuo.mutation.AddSendChunkIDs(ids...)
	return sfuo
}

// AddSendChunks adds the "send_chunks" edges to the SendChunk entity.
func (sfuo *SendFileUpdateOne) AddSendChunks(s ...*SendChunk) *SendFileUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sfuo.AddSendChunkIDs(ids...)
}

// Mutation returns the SendFileMutation object of the builder.
func (sfuo *SendFileUpdateOne) Mutation() *SendFileMutation {
	return sfuo.mutation
}

// ClearSendChunks clears all "send_chunks" edges to the SendChunk entity.
func (sfuo *SendFileUpdateOne) ClearSendChunks() *SendFileUpdateOne {
	sfuo.mutation.ClearSendChunks()
	return sfuo
}

// RemoveSendChunkIDs removes the "send_chunks" edge to SendChunk entities by IDs.
func (sfuo *SendFileUpdateOne) RemoveSendChunkIDs(ids ...int) *SendFileUpdateOne {
	sfuo.mutation.RemoveSendChunkIDs(ids...)
	return sfuo
}

// RemoveSendChunks removes "send_chunks" edges to SendChunk entities.
func (sfuo *SendFileUpdateOne) RemoveSendChunks(s ...*SendChunk) *SendFileUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sfuo.RemoveSendChunkIDs(ids...)
}

// Where appends a list predicates to the SendFileUpdate builder.
func (sfuo *SendFileUpdateOne) Where(ps ...predicate.SendFile) *SendFileUpdateOne {
	sfuo.mutation.Where(ps...)
	return sfuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (sfuo *SendFileUpdateOne) Select(field string, fields ...string) *SendFileUpdateOne {
	sfuo.fields = append([]string{field}, fields...)
	return sfuo
}

// Save executes the query and returns the updated SendFile entity.
func (sfuo *SendFileUpdateOne) Save(ctx context.Context) (*SendFile, error) {
	sfuo.defaults()
	return withHooks(ctx, sfuo.sqlSave, sfuo.mutation, sfuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (sfuo *SendFileUpdateOne) SaveX(ctx context.Context) *SendFile {
	node, err := sfuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (sfuo *SendFileUpdateOne) Exec(ctx context.Context) error {
	_, err := sfuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sfuo *SendFileUpdateOne) ExecX(ctx context.Context) {
	if err := sfuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sfuo *SendFileUpdateOne) defaults() {
	if _, ok := sfuo.mutation.UpdatedAt(); !ok {
		v := sendfile.UpdateDefaultUpdatedAt()
		sfuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sfuo *SendFileUpdateOne) check() error {
	if v, ok := sfuo.mutation.TaskID(); ok {
		if err := sendfile.TaskIDValidator(v); err != nil {
			return &ValidationError{Name: "task_id", err: fmt.Errorf(`ent: validator failed for field "SendFile.task_id": %w`, err)}
		}
	}
	if v, ok := sfuo.mutation.TaskName(); ok {
		if err := sendfile.TaskNameValidator(v); err != nil {
			return &ValidationError{Name: "task_name", err: fmt.Errorf(`ent: validator failed for field "SendFile.task_name": %w`, err)}
		}
	}
	if v, ok := sfuo.mutation.FilePath(); ok {
		if err := sendfile.FilePathValidator(v); err != nil {
			return &ValidationError{Name: "file_path", err: fmt.Errorf(`ent: validator failed for field "SendFile.file_path": %w`, err)}
		}
	}
	if v, ok := sfuo.mutation.Fid(); ok {
		if err := sendfile.FidValidator(v); err != nil {
			return &ValidationError{Name: "fid", err: fmt.Errorf(`ent: validator failed for field "SendFile.fid": %w`, err)}
		}
	}
	if v, ok := sfuo.mutation.Elapsed(); ok {
		if err := sendfile.ElapsedValidator(v); err != nil {
			return &ValidationError{Name: "elapsed", err: fmt.Errorf(`ent: validator failed for field "SendFile.elapsed": %w`, err)}
		}
	}
	if v, ok := sfuo.mutation.Speed(); ok {
		if err := sendfile.SpeedValidator(v); err != nil {
			return &ValidationError{Name: "speed", err: fmt.Errorf(`ent: validator failed for field "SendFile.speed": %w`, err)}
		}
	}
	return nil
}

func (sfuo *SendFileUpdateOne) sqlSave(ctx context.Context) (_node *SendFile, err error) {
	if err := sfuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(sendfile.Table, sendfile.Columns, sqlgraph.NewFieldSpec(sendfile.FieldID, field.TypeInt))
	id, ok := sfuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SendFile.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := sfuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, sendfile.FieldID)
		for _, f := range fields {
			if !sendfile.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != sendfile.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := sfuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sfuo.mutation.TaskID(); ok {
		_spec.SetField(sendfile.FieldTaskID, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.TaskName(); ok {
		_spec.SetField(sendfile.FieldTaskName, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.FilePath(); ok {
		_spec.SetField(sendfile.FieldFilePath, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.Fid(); ok {
		_spec.SetField(sendfile.FieldFid, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.FileSize(); ok {
		_spec.SetField(sendfile.FieldFileSize, field.TypeInt64, value)
	}
	if value, ok := sfuo.mutation.AddedFileSize(); ok {
		_spec.AddField(sendfile.FieldFileSize, field.TypeInt64, value)
	}
	if value, ok := sfuo.mutation.ChunkNumTotal(); ok {
		_spec.SetField(sendfile.FieldChunkNumTotal, field.TypeInt, value)
	}
	if value, ok := sfuo.mutation.AddedChunkNumTotal(); ok {
		_spec.AddField(sendfile.FieldChunkNumTotal, field.TypeInt, value)
	}
	if value, ok := sfuo.mutation.ChunkNumSended(); ok {
		_spec.SetField(sendfile.FieldChunkNumSended, field.TypeInt, value)
	}
	if value, ok := sfuo.mutation.AddedChunkNumSended(); ok {
		_spec.AddField(sendfile.FieldChunkNumSended, field.TypeInt, value)
	}
	if value, ok := sfuo.mutation.Status(); ok {
		_spec.SetField(sendfile.FieldStatus, field.TypeInt, value)
	}
	if value, ok := sfuo.mutation.AddedStatus(); ok {
		_spec.AddField(sendfile.FieldStatus, field.TypeInt, value)
	}
	if value, ok := sfuo.mutation.Elapsed(); ok {
		_spec.SetField(sendfile.FieldElapsed, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.Speed(); ok {
		_spec.SetField(sendfile.FieldSpeed, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.UpdatedAt(); ok {
		_spec.SetField(sendfile.FieldUpdatedAt, field.TypeTime, value)
	}
	if sfuo.mutation.SendChunksCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := sfuo.mutation.RemovedSendChunksIDs(); len(nodes) > 0 && !sfuo.mutation.SendChunksCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := sfuo.mutation.SendChunksIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &SendFile{config: sfuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, sfuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{sendfile.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	sfuo.mutation.done = true
	return _node, nil
}
