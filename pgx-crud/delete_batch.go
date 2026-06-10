package pgx_crud

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// DeleteBatch удаление списка элементов по идентификаторам
func (c *Crud[T]) DeleteBatch(ctx context.Context, ids []int64) error {

	sql, params, err := c.builder.
		Delete(c.table).
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return fmt.Errorf("buildQuery: %w", err)
	}

	_, err = c.db.Exec(ctx, sql, params...)
	if err != nil {
		return fmt.Errorf("execx: %w", err)
	}
	return nil
}
