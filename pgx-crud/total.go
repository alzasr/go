package pgx_crud

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (c *Crud[T]) Total(ctx context.Context, filter Filter) (uint64, error) {

	query := c.builder.Select("count(1)").From(c.table)
	if filter != nil {
		query = filter.ModifyQuery(query)
	}
	query = query.RemoveLimit().RemoveOffset()
	rows, err := c.query(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	res, err := pgx.CollectOneRow(rows, pgx.RowTo[uint64])
	if err != nil {
		return 0, fmt.Errorf("pgx.CollectOneRow: %w", err)
	}
	return res, nil
}
