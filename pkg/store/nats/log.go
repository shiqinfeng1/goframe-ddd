package nats

import (
	"fmt"
)

const (
	uuidLength = 36
)

type Log struct {
	Type     string `json:"type"`
	Duration int64  `json:"duration"`
	Key      string `json:"key"`
	Value    string `json:"value,omitempty"`
}

func (l *Log) String() string {
	var description string

	switch l.Type {
	case "GET":
		description = fmt.Sprintf("Fetching record from bucket '%s' with ID '%s'", l.Value, l.Key)
	case "SET":
		if len(l.Key) == uuidLength {
			description = fmt.Sprintf("Creating new record in bucket '%s' with ID '%s'", l.Value, l.Key)
		} else {
			description = fmt.Sprintf("Updating record with ID '%s' in bucket '%s'", l.Key, l.Value)
		}
	case "DELETE":
		description = fmt.Sprintf("Deleting record from bucket '%s' with ID '%s'", l.Value, l.Key)
	}

	return fmt.Sprintf("%-32s \u001B[38;5;162mNATS\u001B[0m   %8dμs \u001B[38;5;8m%s\u001B[0m\n",
		l.Type,
		l.Duration,
		description)
}
