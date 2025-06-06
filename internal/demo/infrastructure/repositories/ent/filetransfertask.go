// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/shiqinfeng1/goframe-ddd/internal/demo/infrastructure/repositories/ent/filetransfertask"
)

// FileTransferTask is the model entity for the FileTransferTask schema.
type FileTransferTask struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// TaskID holds the value of the "task_id" field.
	TaskID string `json:"task_id,omitempty"`
	// TaskName holds the value of the "task_name" field.
	TaskName string `json:"task_name,omitempty"`
	// NodeID holds the value of the "node_id" field.
	NodeID string `json:"node_id,omitempty"`
	// Status holds the value of the "status" field.
	Status int `json:"status,omitempty"`
	// Elapsed holds the value of the "elapsed" field.
	Elapsed string `json:"elapsed,omitempty"`
	// Speed holds the value of the "speed" field.
	Speed string `json:"speed,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt    time.Time `json:"created_at,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*FileTransferTask) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case filetransfertask.FieldID, filetransfertask.FieldStatus:
			values[i] = new(sql.NullInt64)
		case filetransfertask.FieldTaskID, filetransfertask.FieldTaskName, filetransfertask.FieldNodeID, filetransfertask.FieldElapsed, filetransfertask.FieldSpeed:
			values[i] = new(sql.NullString)
		case filetransfertask.FieldUpdatedAt, filetransfertask.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the FileTransferTask fields.
func (ftt *FileTransferTask) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case filetransfertask.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ftt.ID = int(value.Int64)
		case filetransfertask.FieldTaskID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field task_id", values[i])
			} else if value.Valid {
				ftt.TaskID = value.String
			}
		case filetransfertask.FieldTaskName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field task_name", values[i])
			} else if value.Valid {
				ftt.TaskName = value.String
			}
		case filetransfertask.FieldNodeID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field node_id", values[i])
			} else if value.Valid {
				ftt.NodeID = value.String
			}
		case filetransfertask.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				ftt.Status = int(value.Int64)
			}
		case filetransfertask.FieldElapsed:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field elapsed", values[i])
			} else if value.Valid {
				ftt.Elapsed = value.String
			}
		case filetransfertask.FieldSpeed:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field speed", values[i])
			} else if value.Valid {
				ftt.Speed = value.String
			}
		case filetransfertask.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ftt.UpdatedAt = value.Time
			}
		case filetransfertask.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ftt.CreatedAt = value.Time
			}
		default:
			ftt.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the FileTransferTask.
// This includes values selected through modifiers, order, etc.
func (ftt *FileTransferTask) Value(name string) (ent.Value, error) {
	return ftt.selectValues.Get(name)
}

// Update returns a builder for updating this FileTransferTask.
// Note that you need to call FileTransferTask.Unwrap() before calling this method if this FileTransferTask
// was returned from a transaction, and the transaction was committed or rolled back.
func (ftt *FileTransferTask) Update() *FileTransferTaskUpdateOne {
	return NewFileTransferTaskClient(ftt.config).UpdateOne(ftt)
}

// Unwrap unwraps the FileTransferTask entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ftt *FileTransferTask) Unwrap() *FileTransferTask {
	_tx, ok := ftt.config.driver.(*txDriver)
	if !ok {
		panic("ent: FileTransferTask is not a transactional entity")
	}
	ftt.config.driver = _tx.drv
	return ftt
}

// String implements the fmt.Stringer.
func (ftt *FileTransferTask) String() string {
	var builder strings.Builder
	builder.WriteString("FileTransferTask(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ftt.ID))
	builder.WriteString("task_id=")
	builder.WriteString(ftt.TaskID)
	builder.WriteString(", ")
	builder.WriteString("task_name=")
	builder.WriteString(ftt.TaskName)
	builder.WriteString(", ")
	builder.WriteString("node_id=")
	builder.WriteString(ftt.NodeID)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", ftt.Status))
	builder.WriteString(", ")
	builder.WriteString("elapsed=")
	builder.WriteString(ftt.Elapsed)
	builder.WriteString(", ")
	builder.WriteString("speed=")
	builder.WriteString(ftt.Speed)
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(ftt.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(ftt.CreatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// FileTransferTasks is a parsable slice of FileTransferTask.
type FileTransferTasks []*FileTransferTask
