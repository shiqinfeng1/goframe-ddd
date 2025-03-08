// Code generated by ent, DO NOT EDIT.

package sendfile

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldID, id))
}

// TaskID applies equality check predicate on the "task_id" field. It's identical to TaskIDEQ.
func TaskID(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldTaskID, v))
}

// TaskName applies equality check predicate on the "task_name" field. It's identical to TaskNameEQ.
func TaskName(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldTaskName, v))
}

// FilePath applies equality check predicate on the "file_path" field. It's identical to FilePathEQ.
func FilePath(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldFilePath, v))
}

// FileID applies equality check predicate on the "file_id" field. It's identical to FileIDEQ.
func FileID(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldFileID, v))
}

// FileSize applies equality check predicate on the "file_size" field. It's identical to FileSizeEQ.
func FileSize(v int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldFileSize, v))
}

// ChunkNumTotal applies equality check predicate on the "chunk_num_total" field. It's identical to ChunkNumTotalEQ.
func ChunkNumTotal(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldChunkNumTotal, v))
}

// ChunkNumSended applies equality check predicate on the "chunk_num_sended" field. It's identical to ChunkNumSendedEQ.
func ChunkNumSended(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldChunkNumSended, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldStatus, v))
}

// Elapsed applies equality check predicate on the "elapsed" field. It's identical to ElapsedEQ.
func Elapsed(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldElapsed, v))
}

// Speed applies equality check predicate on the "speed" field. It's identical to SpeedEQ.
func Speed(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldSpeed, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldUpdatedAt, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldCreatedAt, v))
}

// TaskIDEQ applies the EQ predicate on the "task_id" field.
func TaskIDEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldTaskID, v))
}

// TaskIDNEQ applies the NEQ predicate on the "task_id" field.
func TaskIDNEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldTaskID, v))
}

// TaskIDIn applies the In predicate on the "task_id" field.
func TaskIDIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldTaskID, vs...))
}

// TaskIDNotIn applies the NotIn predicate on the "task_id" field.
func TaskIDNotIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldTaskID, vs...))
}

// TaskIDGT applies the GT predicate on the "task_id" field.
func TaskIDGT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldTaskID, v))
}

// TaskIDGTE applies the GTE predicate on the "task_id" field.
func TaskIDGTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldTaskID, v))
}

// TaskIDLT applies the LT predicate on the "task_id" field.
func TaskIDLT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldTaskID, v))
}

// TaskIDLTE applies the LTE predicate on the "task_id" field.
func TaskIDLTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldTaskID, v))
}

// TaskIDContains applies the Contains predicate on the "task_id" field.
func TaskIDContains(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContains(FieldTaskID, v))
}

// TaskIDHasPrefix applies the HasPrefix predicate on the "task_id" field.
func TaskIDHasPrefix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasPrefix(FieldTaskID, v))
}

// TaskIDHasSuffix applies the HasSuffix predicate on the "task_id" field.
func TaskIDHasSuffix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasSuffix(FieldTaskID, v))
}

// TaskIDEqualFold applies the EqualFold predicate on the "task_id" field.
func TaskIDEqualFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEqualFold(FieldTaskID, v))
}

// TaskIDContainsFold applies the ContainsFold predicate on the "task_id" field.
func TaskIDContainsFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContainsFold(FieldTaskID, v))
}

// TaskNameEQ applies the EQ predicate on the "task_name" field.
func TaskNameEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldTaskName, v))
}

// TaskNameNEQ applies the NEQ predicate on the "task_name" field.
func TaskNameNEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldTaskName, v))
}

// TaskNameIn applies the In predicate on the "task_name" field.
func TaskNameIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldTaskName, vs...))
}

// TaskNameNotIn applies the NotIn predicate on the "task_name" field.
func TaskNameNotIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldTaskName, vs...))
}

// TaskNameGT applies the GT predicate on the "task_name" field.
func TaskNameGT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldTaskName, v))
}

// TaskNameGTE applies the GTE predicate on the "task_name" field.
func TaskNameGTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldTaskName, v))
}

// TaskNameLT applies the LT predicate on the "task_name" field.
func TaskNameLT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldTaskName, v))
}

// TaskNameLTE applies the LTE predicate on the "task_name" field.
func TaskNameLTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldTaskName, v))
}

// TaskNameContains applies the Contains predicate on the "task_name" field.
func TaskNameContains(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContains(FieldTaskName, v))
}

// TaskNameHasPrefix applies the HasPrefix predicate on the "task_name" field.
func TaskNameHasPrefix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasPrefix(FieldTaskName, v))
}

// TaskNameHasSuffix applies the HasSuffix predicate on the "task_name" field.
func TaskNameHasSuffix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasSuffix(FieldTaskName, v))
}

// TaskNameEqualFold applies the EqualFold predicate on the "task_name" field.
func TaskNameEqualFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEqualFold(FieldTaskName, v))
}

// TaskNameContainsFold applies the ContainsFold predicate on the "task_name" field.
func TaskNameContainsFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContainsFold(FieldTaskName, v))
}

// FilePathEQ applies the EQ predicate on the "file_path" field.
func FilePathEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldFilePath, v))
}

// FilePathNEQ applies the NEQ predicate on the "file_path" field.
func FilePathNEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldFilePath, v))
}

// FilePathIn applies the In predicate on the "file_path" field.
func FilePathIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldFilePath, vs...))
}

// FilePathNotIn applies the NotIn predicate on the "file_path" field.
func FilePathNotIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldFilePath, vs...))
}

// FilePathGT applies the GT predicate on the "file_path" field.
func FilePathGT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldFilePath, v))
}

// FilePathGTE applies the GTE predicate on the "file_path" field.
func FilePathGTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldFilePath, v))
}

// FilePathLT applies the LT predicate on the "file_path" field.
func FilePathLT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldFilePath, v))
}

// FilePathLTE applies the LTE predicate on the "file_path" field.
func FilePathLTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldFilePath, v))
}

// FilePathContains applies the Contains predicate on the "file_path" field.
func FilePathContains(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContains(FieldFilePath, v))
}

// FilePathHasPrefix applies the HasPrefix predicate on the "file_path" field.
func FilePathHasPrefix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasPrefix(FieldFilePath, v))
}

// FilePathHasSuffix applies the HasSuffix predicate on the "file_path" field.
func FilePathHasSuffix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasSuffix(FieldFilePath, v))
}

// FilePathEqualFold applies the EqualFold predicate on the "file_path" field.
func FilePathEqualFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEqualFold(FieldFilePath, v))
}

// FilePathContainsFold applies the ContainsFold predicate on the "file_path" field.
func FilePathContainsFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContainsFold(FieldFilePath, v))
}

// FileIDEQ applies the EQ predicate on the "file_id" field.
func FileIDEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldFileID, v))
}

// FileIDNEQ applies the NEQ predicate on the "file_id" field.
func FileIDNEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldFileID, v))
}

// FileIDIn applies the In predicate on the "file_id" field.
func FileIDIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldFileID, vs...))
}

// FileIDNotIn applies the NotIn predicate on the "file_id" field.
func FileIDNotIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldFileID, vs...))
}

// FileIDGT applies the GT predicate on the "file_id" field.
func FileIDGT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldFileID, v))
}

// FileIDGTE applies the GTE predicate on the "file_id" field.
func FileIDGTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldFileID, v))
}

// FileIDLT applies the LT predicate on the "file_id" field.
func FileIDLT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldFileID, v))
}

// FileIDLTE applies the LTE predicate on the "file_id" field.
func FileIDLTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldFileID, v))
}

// FileIDContains applies the Contains predicate on the "file_id" field.
func FileIDContains(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContains(FieldFileID, v))
}

// FileIDHasPrefix applies the HasPrefix predicate on the "file_id" field.
func FileIDHasPrefix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasPrefix(FieldFileID, v))
}

// FileIDHasSuffix applies the HasSuffix predicate on the "file_id" field.
func FileIDHasSuffix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasSuffix(FieldFileID, v))
}

// FileIDEqualFold applies the EqualFold predicate on the "file_id" field.
func FileIDEqualFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEqualFold(FieldFileID, v))
}

// FileIDContainsFold applies the ContainsFold predicate on the "file_id" field.
func FileIDContainsFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContainsFold(FieldFileID, v))
}

// FileSizeEQ applies the EQ predicate on the "file_size" field.
func FileSizeEQ(v int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldFileSize, v))
}

// FileSizeNEQ applies the NEQ predicate on the "file_size" field.
func FileSizeNEQ(v int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldFileSize, v))
}

// FileSizeIn applies the In predicate on the "file_size" field.
func FileSizeIn(vs ...int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldFileSize, vs...))
}

// FileSizeNotIn applies the NotIn predicate on the "file_size" field.
func FileSizeNotIn(vs ...int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldFileSize, vs...))
}

// FileSizeGT applies the GT predicate on the "file_size" field.
func FileSizeGT(v int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldFileSize, v))
}

// FileSizeGTE applies the GTE predicate on the "file_size" field.
func FileSizeGTE(v int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldFileSize, v))
}

// FileSizeLT applies the LT predicate on the "file_size" field.
func FileSizeLT(v int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldFileSize, v))
}

// FileSizeLTE applies the LTE predicate on the "file_size" field.
func FileSizeLTE(v int64) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldFileSize, v))
}

// ChunkNumTotalEQ applies the EQ predicate on the "chunk_num_total" field.
func ChunkNumTotalEQ(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldChunkNumTotal, v))
}

// ChunkNumTotalNEQ applies the NEQ predicate on the "chunk_num_total" field.
func ChunkNumTotalNEQ(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldChunkNumTotal, v))
}

// ChunkNumTotalIn applies the In predicate on the "chunk_num_total" field.
func ChunkNumTotalIn(vs ...int) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldChunkNumTotal, vs...))
}

// ChunkNumTotalNotIn applies the NotIn predicate on the "chunk_num_total" field.
func ChunkNumTotalNotIn(vs ...int) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldChunkNumTotal, vs...))
}

// ChunkNumTotalGT applies the GT predicate on the "chunk_num_total" field.
func ChunkNumTotalGT(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldChunkNumTotal, v))
}

// ChunkNumTotalGTE applies the GTE predicate on the "chunk_num_total" field.
func ChunkNumTotalGTE(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldChunkNumTotal, v))
}

// ChunkNumTotalLT applies the LT predicate on the "chunk_num_total" field.
func ChunkNumTotalLT(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldChunkNumTotal, v))
}

// ChunkNumTotalLTE applies the LTE predicate on the "chunk_num_total" field.
func ChunkNumTotalLTE(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldChunkNumTotal, v))
}

// ChunkNumSendedEQ applies the EQ predicate on the "chunk_num_sended" field.
func ChunkNumSendedEQ(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldChunkNumSended, v))
}

// ChunkNumSendedNEQ applies the NEQ predicate on the "chunk_num_sended" field.
func ChunkNumSendedNEQ(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldChunkNumSended, v))
}

// ChunkNumSendedIn applies the In predicate on the "chunk_num_sended" field.
func ChunkNumSendedIn(vs ...int) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldChunkNumSended, vs...))
}

// ChunkNumSendedNotIn applies the NotIn predicate on the "chunk_num_sended" field.
func ChunkNumSendedNotIn(vs ...int) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldChunkNumSended, vs...))
}

// ChunkNumSendedGT applies the GT predicate on the "chunk_num_sended" field.
func ChunkNumSendedGT(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldChunkNumSended, v))
}

// ChunkNumSendedGTE applies the GTE predicate on the "chunk_num_sended" field.
func ChunkNumSendedGTE(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldChunkNumSended, v))
}

// ChunkNumSendedLT applies the LT predicate on the "chunk_num_sended" field.
func ChunkNumSendedLT(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldChunkNumSended, v))
}

// ChunkNumSendedLTE applies the LTE predicate on the "chunk_num_sended" field.
func ChunkNumSendedLTE(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldChunkNumSended, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...int) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...int) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v int) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldStatus, v))
}

// ElapsedEQ applies the EQ predicate on the "elapsed" field.
func ElapsedEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldElapsed, v))
}

// ElapsedNEQ applies the NEQ predicate on the "elapsed" field.
func ElapsedNEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldElapsed, v))
}

// ElapsedIn applies the In predicate on the "elapsed" field.
func ElapsedIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldElapsed, vs...))
}

// ElapsedNotIn applies the NotIn predicate on the "elapsed" field.
func ElapsedNotIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldElapsed, vs...))
}

// ElapsedGT applies the GT predicate on the "elapsed" field.
func ElapsedGT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldElapsed, v))
}

// ElapsedGTE applies the GTE predicate on the "elapsed" field.
func ElapsedGTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldElapsed, v))
}

// ElapsedLT applies the LT predicate on the "elapsed" field.
func ElapsedLT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldElapsed, v))
}

// ElapsedLTE applies the LTE predicate on the "elapsed" field.
func ElapsedLTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldElapsed, v))
}

// ElapsedContains applies the Contains predicate on the "elapsed" field.
func ElapsedContains(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContains(FieldElapsed, v))
}

// ElapsedHasPrefix applies the HasPrefix predicate on the "elapsed" field.
func ElapsedHasPrefix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasPrefix(FieldElapsed, v))
}

// ElapsedHasSuffix applies the HasSuffix predicate on the "elapsed" field.
func ElapsedHasSuffix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasSuffix(FieldElapsed, v))
}

// ElapsedEqualFold applies the EqualFold predicate on the "elapsed" field.
func ElapsedEqualFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEqualFold(FieldElapsed, v))
}

// ElapsedContainsFold applies the ContainsFold predicate on the "elapsed" field.
func ElapsedContainsFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContainsFold(FieldElapsed, v))
}

// SpeedEQ applies the EQ predicate on the "speed" field.
func SpeedEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldSpeed, v))
}

// SpeedNEQ applies the NEQ predicate on the "speed" field.
func SpeedNEQ(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldSpeed, v))
}

// SpeedIn applies the In predicate on the "speed" field.
func SpeedIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldSpeed, vs...))
}

// SpeedNotIn applies the NotIn predicate on the "speed" field.
func SpeedNotIn(vs ...string) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldSpeed, vs...))
}

// SpeedGT applies the GT predicate on the "speed" field.
func SpeedGT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldSpeed, v))
}

// SpeedGTE applies the GTE predicate on the "speed" field.
func SpeedGTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldSpeed, v))
}

// SpeedLT applies the LT predicate on the "speed" field.
func SpeedLT(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldSpeed, v))
}

// SpeedLTE applies the LTE predicate on the "speed" field.
func SpeedLTE(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldSpeed, v))
}

// SpeedContains applies the Contains predicate on the "speed" field.
func SpeedContains(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContains(FieldSpeed, v))
}

// SpeedHasPrefix applies the HasPrefix predicate on the "speed" field.
func SpeedHasPrefix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasPrefix(FieldSpeed, v))
}

// SpeedHasSuffix applies the HasSuffix predicate on the "speed" field.
func SpeedHasSuffix(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldHasSuffix(FieldSpeed, v))
}

// SpeedEqualFold applies the EqualFold predicate on the "speed" field.
func SpeedEqualFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldEqualFold(FieldSpeed, v))
}

// SpeedContainsFold applies the ContainsFold predicate on the "speed" field.
func SpeedContainsFold(v string) predicate.SendFile {
	return predicate.SendFile(sql.FieldContainsFold(FieldSpeed, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldUpdatedAt, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.SendFile {
	return predicate.SendFile(sql.FieldLTE(FieldCreatedAt, v))
}

// HasSendChunks applies the HasEdge predicate on the "send_chunks" edge.
func HasSendChunks() predicate.SendFile {
	return predicate.SendFile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SendChunksTable, SendChunksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSendChunksWith applies the HasEdge predicate on the "send_chunks" edge with a given conditions (other predicates).
func HasSendChunksWith(preds ...predicate.SendChunk) predicate.SendFile {
	return predicate.SendFile(func(s *sql.Selector) {
		step := newSendChunksStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.SendFile) predicate.SendFile {
	return predicate.SendFile(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.SendFile) predicate.SendFile {
	return predicate.SendFile(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.SendFile) predicate.SendFile {
	return predicate.SendFile(sql.NotPredicates(p))
}
