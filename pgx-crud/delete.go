package pgx_crud

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// Delete удалить элемент по id
func (c *Crud[T]) Delete(ctx context.Context, id int64) (int64, error) {

	query := c.builder.
		Delete(c.table).
		Where(squirrel.Eq{"id": id})

	rowsAffected, err := c.exec(ctx, query)
	if err != nil {
		return rowsAffected, fmt.Errorf("exec: %w", err)
	}
	
	return rowsAffected, nil
}
