package pgx_crud

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

type DeleteFilter interface {
	ModifyQuery(builder squirrel.DeleteBuilder) squirrel.DeleteBuilder
}

func (c *Crud[T]) DeleteByFilter(ctx context.Context, filter DeleteFilter) (affectedRows int64, err error) {

	query := c.builder.
		Delete(c.table)

	query = filter.ModifyQuery(query)

	affectedRows, err = c.exec(ctx, query)

	if err != nil {
		return affectedRows, fmt.Errorf("exec: %w", err)
	}
	return affectedRows, nil
}
