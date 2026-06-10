package pgx_crud

import (
	"context"
	"fmt"
)

// CreateBatch создать несколько элементов одним запросом
func (c *Crud[T]) CreateBatch(ctx context.Context, models []*T) ([]*T, error) {
	if len(models) == 0 {
		return nil, nil
	}
	query := c.builder.
		Insert(c.table).
		Columns(c.insertFields...).
		Suffix(c.returning())

	for i := range models {
		query = query.Values(c.valuesSlice(models[i], c.insertFields)...)
	}

	res, err := c.queryAll(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("queryAll: %w", err)
	}

	return res, nil
}
