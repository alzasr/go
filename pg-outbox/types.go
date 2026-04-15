package pgoutbox

import (
	"context"
	"database/sql"
	"errors"
)

type Tx interface {
	Query(ctx context.Context, sql string, params ...any) (*sql.Rows, error)
}

type MsgStore interface {
	SaveTx(ctx context.Context, tx Tx, msg *Message) error
}

type Logger interface {
	Error(err error)
}

type Locker interface {
	TryLock(ctx context.Context, key []byte) (locked bool, unlock func() error, err error)
}

type Sender interface {
	Send(ctx context.Context, msgs []*Message) error
}

var noop = func() error { return nil }

var ErrNotFound = errors.New("not found")
