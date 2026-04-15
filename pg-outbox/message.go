package pgoutbox

import "time"

type Message struct {
	// SYSTEM FIELDS. DO NOT FILL
	ID        int64
	Attempt   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	LastError string

	// User Fields
	Topic   string
	Key     []byte
	Headers []byte
	Payload []byte
}
