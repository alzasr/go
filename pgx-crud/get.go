package pgx_crud

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

// Get получение элемента по идентификатору
func (c *Crud[T]) Get(ctx context.Context, id int64) (*T, error) {
	query := c.builder.
		Select(c.selectFields...).
		From(c.table).
		Where(squirrel.Eq{"id": id})
	
	res, err := c.queryOne(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, c.getNotFoundErrorOrErr(fmt.Errorf("getx: %w", err))
		}
		return nil, fmt.Errorf("getx: %w", err)
	}
	return res, nil
}
