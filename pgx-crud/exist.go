package pgx_crud

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (c *Crud[T]) Exist(ctx context.Context, filter Filter) (bool, error) {
	query := c.builder.Select("1").From(c.table)
	if filter != nil {
		query = filter.ModifyQuery(query)
	}
	query = query.RemoveOffset().Limit(1)
	rows, err := c.query(ctx, query)
	if err != nil {
		return false, fmt.Errorf("query: %w", err)
	}
	_ , err = pgx.CollectOneRow(rows, pgx.RowTo[int8] )
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, pgx.ErrNoRows):
		return false, nil
	default:
		return false, fmt.Errorf("pgx.CollectOneRow: %w", err)
	}
}
