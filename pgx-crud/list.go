package pgx_crud

import (
	"context"
	"fmt"
	
)

// List получение списка элементов согласно фильтру
func (c *Crud[T]) List(ctx context.Context, filter Filter) ([]*T, error) {
	query := c.builder.Select(c.selectFields...).From(c.table)
	if filter != nil {
		query = filter.ModifyQuery(query)
	}
	res, err := c.queryAll(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("queryAll: %w", err)
	}
	return res, nil
}
