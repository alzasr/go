package pgoutbox

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type store struct {
	builder sq.StatementBuilderType
	table   string
	db      Tx
}

func NewPGStore(db Tx) *store {
	return &store{
		sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		"outbox",
		db,
	}
}

func (s *store) CreateTx(ctx context.Context, tx Tx, msg *Message) error {
	now := time.Now()
	msg.CreatedAt = now
	msg.UpdatedAt = now
	msg.LastError = ""
	msg.Attempt = 0

	sql, params, err := s.builder.
		Insert(s.table).
		SetMap(map[string]interface{}{
			"attempt":    msg.Attempt,
			"created_at": msg.CreatedAt,
			"updated_at": msg.UpdatedAt,
			"last_error": msg.LastError,
			"topic":      msg.Topic,
			"key":        msg.Key,
			"headers":    msg.Headers,
			"payload":    msg.Payload,
		}).Suffix(`RETURNING "id"`).ToSql()
	if err != nil {
		return fmt.Errorf("build create query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, params...)
	if err != nil {
		return fmt.Errorf("query insert: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("undefined error, dont returning id")
	}

	err = rows.Scan(&msg.ID)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	return nil
}

func (s *store) DeleteBatch(ctx context.Context, msgs []*Message) error {
	if len(msgs) == 0 {
		return nil
	}

	ids := make([]int64, len(msgs))
	for i := range msgs {
		ids[i] = msgs[i].ID
	}

	sql, params, err := s.builder.
		Delete(s.table).
		Where(sq.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, params...)
	if err != nil {
		return fmt.Errorf("exec delete query: %w", err)
	}
	defer rows.Close()

	return nil
}

func (s *store) UpdateBatch(ctx context.Context, msgs []*Message) error {
	now := time.Now()
	for i := range msgs {
		msgs[i].UpdatedAt = now
		err := s.Update(ctx, msgs[i])
		if err != nil {
			return fmt.Errorf("update msg(%v): %w", msgs[i].ID, err)
		}
	}
	return nil
}

func (s *store) Update(ctx context.Context, msg *Message) error {
	sql, params, err := s.builder.Update(s.table).SetMap(map[string]interface{}{
		"attempt":    msg.Attempt,
		"updated_at": msg.UpdatedAt,
		"last_error": msg.LastError,
	}).Where(sq.Eq{"id": msg.ID}).ToSql()

	rows, err := s.db.Query(ctx, sql, params...)
	if err != nil {
		return fmt.Errorf("exec query update: %w", err)
	}
	defer rows.Close()
	return nil
}

func (s *store) ListByKey(ctx context.Context, key []byte) ([]*Message, error) {
	sql, params, err := s.builder.
		Select(s.selectFields()...).
		From(s.table).
		Where(sq.Eq{"key": key}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}
	rows, err := s.db.Query(ctx, sql, params...)
	if err != nil {
		return nil, fmt.Errorf("exec select query: %w", err)
	}

	defer rows.Close()

	var msgs []*Message

	for rows.Next() {
		msg, err := s.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan rows: %w", err)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (s *store) ListIDsGroupingByKey(ctx context.Context) ([]int64, error) {
	sql, params, err := s.builder.Select(`min("id")`).From(s.table).OrderBy("updated_at").GroupBy("key").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select min id query: %w", err)
	}
	rows, err := s.db.Query(ctx, sql, params...)
	if err != nil {
		return nil, fmt.Errorf("exec select min id query: %w", err)
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("scan id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *store) Get(ctx context.Context, id int64) (*Message, error) {
	sql, params, err := s.builder.Select(s.selectFields()...).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get query: %w", err)
	}
	rows, err := s.db.Query(ctx, sql, params...)
	if err != nil {
		return nil, fmt.Errorf("exec get query: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, ErrNotFound
	}
	msg, err := s.scan(rows)
	if err != nil {
		return nil, fmt.Errorf("scan rows: %w", err)
	}
	return msg, nil
}

func (s *store) selectFields() []string {
	return []string{"id", "attempt", "created_at", "updated_at", "last_error", "topic", "key", "headers", "payload"}
}

func (s *store) scan(rows *sql.Rows) (*Message, error) {
	msg := Message{}
	err := rows.Scan(&msg.ID, &msg.Attempt, &msg.CreatedAt, &msg.UpdatedAt, &msg.LastError, &msg.Topic, &msg.Key, &msg.Headers, &msg.Payload)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return &msg, nil
}
